package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/flags/internal/server"
	"github.com/w-h-a/flags/internal/server/clients/exporter"
	localexporter "github.com/w-h-a/flags/internal/server/clients/exporter/local"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
	localnotifier "github.com/w-h-a/flags/internal/server/clients/notifier/local"
	"github.com/w-h-a/flags/internal/server/clients/notifier/slack"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/clients/reader/github"
	"github.com/w-h-a/flags/internal/server/clients/reader/gitlab"
	localreader "github.com/w-h-a/flags/internal/server/clients/reader/local"
	postgresreader "github.com/w-h-a/flags/internal/server/clients/reader/postgres"
	"github.com/w-h-a/flags/internal/server/clients/writer"
	"github.com/w-h-a/flags/internal/server/clients/writer/noop"
	postgreswriter "github.com/w-h-a/flags/internal/server/clients/writer/postgres"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Server(ctx *cli.Context) error {
	// config
	config.New()

	// resource
	name := config.Name()

	instCtx := context.Background()

	resource, err := resource.New(
		instCtx,
		resource.WithAttributes(
			semconv.ServiceName(name),
		),
		resource.WithProcess(),
	)
	if err != nil {
		return err
	}

	// log
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix(name),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

	// traces
	traceExporter, err := otlptracehttp.New(
		instCtx,
		otlptracehttp.WithEndpoint(config.TracesAddress()),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(
			trace.NewBatchSpanProcessor(
				traceExporter,
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	defer func() {
		if err := tp.Shutdown(instCtx); err != nil {
			log.Warnf("failed to gracefully shutdown trace provider: %v", err)
		}
	}()

	// metrics
	metricsExporter, err := otlpmetrichttp.New(
		instCtx,
		otlpmetrichttp.WithEndpoint(config.MetricsAddress()),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricsExporter,
				metric.WithInterval(15*time.Second),
				metric.WithProducer(runtime.NewProducer()),
			),
		),
	)

	otel.SetMeterProvider(mp)
	defer func() {
		if err := mp.Shutdown(instCtx); err != nil {
			log.Warnf("failed to gracefully shutdown metric provider: %v", err)
		}
	}()

	if err := host.Start(); err != nil {
		return err
	}

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return err
	}

	// clients
	writeClient := initWriteClient()
	readClient := initReadClient()
	exportClient := initExportClient()
	notifyClient := initNotifyClient()

	// server + services
	httpServer, cacheService, exportService, notifyService, err := server.Factory(
		writeClient,
		readClient,
		exportClient,
		notifyClient,
	)
	if err != nil {
		return err
	}

	// wait group and error chan
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 3)

	// start http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- httpServer.Start()
	}()

	// start cache updater
	cacheStop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- server.UpdateCache(
			cacheService,
			notifyService,
			cacheStop,
			time.Duration(config.ReadInterval())*time.Second,
		)
	}()

	// start exporter
	exportStop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- server.ExportReports(
			exportService,
			exportStop,
			time.Duration(config.ExportInterval())*time.Second,
		)
	}()

	// block
	err = <-errCh
	if err != nil {
		log.Errorf("failed to start server: %+v", err)
		return err
	}

	// graceful shutdown
	wait := make(chan struct{})

	go func() {
		defer close(wait)
		wg.Wait()
	}()

	close(exportStop)

	log.Info("successfully stopped export")

	close(cacheStop)

	log.Info("successfully stopped cache")

	select {
	case <-wait:
	case <-time.After(30 * time.Second):
	}

	log.Info("successfully stopped server")

	return nil
}

func initWriteClient() writer.Writer {
	switch config.WriteClient() {
	case "postgres":
		return postgreswriter.NewWriter(
			writer.WithLocation(config.WriteClientLocation()),
		)
	default:
		return noop.NewWriter(
			writer.WithLocation(config.WriteClientLocation()),
		)
	}
}

func initReadClient() reader.Reader {
	switch config.ReadClient() {
	case "github":
		return github.NewReader(
			reader.WithLocation(config.ReadClientLocation()),
			reader.WithToken(config.ReadClientToken()),
		)
	case "gitlab":
		return gitlab.NewReader(
			reader.WithLocation(config.ReadClientLocation()),
			reader.WithToken(config.ReadClientToken()),
		)
	case "postgres":
		return postgresreader.NewReader(
			reader.WithLocation(config.ReadClientLocation()),
		)
	default:
		return localreader.NewReader(
			reader.WithLocation(config.ReadClientLocation()),
		)
	}
}

func initExportClient() exporter.Exporter {
	switch config.ExportClient() {
	default:
		return localexporter.NewExporter(
			exporter.WithDir(config.ExportClientDir()),
		)
	}
}

func initNotifyClient() notifier.Notifier {
	switch config.NotifyClient() {
	case "slack":
		return slack.NewNotifier(
			notifier.WithURL(config.NotifyURL()),
		)
	default:
		return localnotifier.NewNotifier()
	}
}
