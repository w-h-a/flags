package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/flags/internal/server"
	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/flags/internal/server/clients/file/github"
	"github.com/w-h-a/flags/internal/server/clients/file/gitlab"
	localfile "github.com/w-h-a/flags/internal/server/clients/file/local"
	"github.com/w-h-a/flags/internal/server/clients/message"
	localmessage "github.com/w-h-a/flags/internal/server/clients/message/local"
	"github.com/w-h-a/flags/internal/server/clients/message/slack"
	"github.com/w-h-a/flags/internal/server/clients/report"
	localreport "github.com/w-h-a/flags/internal/server/clients/report/local"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/export"
	"github.com/w-h-a/flags/internal/server/services/notify"
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
	otel.SetTextMapPropagator(propagation.TraceContext{})
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
	fileClient := initFileClient()
	reportClient := initReportClient()
	messageClient := initMessageClient()

	// server + services
	httpServer, cacheService, exportService, notifyService, err := server.Factory(fileClient, reportClient, messageClient)
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
		errCh <- updateCache(cacheService, notifyService, cacheStop)
	}()

	// start exporter
	exportStop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- exportReports(exportService, exportStop)
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

func updateCache(cacheService *cache.Service, notifyService *notify.Service, stop chan struct{}) error {
	// TODO: confirm poll interval is valid

	// TODO: retrieve from config
	ticker := time.NewTicker(2 * time.Minute)

	for {
		select {
		case <-ticker.C:
			old, new, err := cacheService.RetrieveFlags()
			if err != nil {
				log.Warnf("failed to update the cache: %v", err)
			}

			notifyService.Notify(old, new)
		case <-stop:
			ticker.Stop()
			notifyService.Close()
			return nil
		}
	}
}

func exportReports(exportService *export.Service, stop chan struct{}) error {
	// TODO: confirm

	// TODO: retrieve from config
	ticker := time.NewTicker(2 * time.Minute)

	for {
		select {
		case <-ticker.C:
			exportService.Flush()
		case <-stop:
			ticker.Stop()
			exportService.Close()
			return nil
		}
	}
}

func initFileClient() file.Client {
	switch config.FileClient() {
	case "github":
		return github.NewFileClient(
			file.WithDir(config.FileClientDir()),
			file.WithFiles(config.FileClientFiles()...),
			file.WithToken(config.FileClientToken()),
		)
	case "gitlab":
		return gitlab.NewFileClient(
			file.WithDir(config.FileClientDir()),
			file.WithFiles(config.FileClientFiles()...),
			file.WithToken(config.FileClientToken()),
		)
	default:
		return localfile.NewFileClient(
			file.WithDir(config.FileClientDir()),
			file.WithFiles(config.FileClientFiles()...),
		)
	}
}

func initReportClient() report.Client {
	switch config.ReportClient() {
	default:
		return localreport.NewReportClient(
			report.WithDir(config.ReportClientDir()),
		)
	}
}

func initMessageClient() message.Client {
	switch config.MessageClient() {
	case "slack":
		return slack.NewMessageClient(
			message.WithURL(config.MessageURL()),
		)
	default:
		return localmessage.NewMessageClient()
	}
}
