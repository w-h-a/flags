package cmd

import (
	"context"
	"fmt"
	"log/slog"
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
	dynamodbreader "github.com/w-h-a/flags/internal/server/clients/reader/dynamodb"
	"github.com/w-h-a/flags/internal/server/clients/reader/github"
	"github.com/w-h-a/flags/internal/server/clients/reader/gitlab"
	localreader "github.com/w-h-a/flags/internal/server/clients/reader/local"
	postgresreader "github.com/w-h-a/flags/internal/server/clients/reader/postgres"
	"github.com/w-h-a/flags/internal/server/clients/writer"
	dynamodbwriter "github.com/w-h-a/flags/internal/server/clients/writer/dynamodb"
	"github.com/w-h-a/flags/internal/server/clients/writer/noop"
	postgreswriter "github.com/w-h-a/flags/internal/server/clients/writer/postgres"
	"github.com/w-h-a/flags/internal/server/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	globallog "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	logsdk "go.opentelemetry.io/otel/sdk/log"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Server(ctx *cli.Context) error {
	// config
	config.New()

	// resource
	name := config.Name()

	resource, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(name),
		),
		resource.WithProcess(),
	)
	if err != nil {
		return err
	}

	// logs
	logsExporter, err := initLogsExporter(context.Background())
	if err != nil {
		return err
	}

	lp := logsdk.NewLoggerProvider(
		logsdk.WithResource(resource),
		logsdk.WithProcessor(
			logsdk.NewBatchProcessor(logsExporter),
		),
	)

	globallog.SetLoggerProvider(lp)

	defer lp.Shutdown(context.Background())

	logger := otelslog.NewLogger(
		config.Name(),
		otelslog.WithLoggerProvider(lp),
	)

	slog.SetDefault(logger)

	// traces
	traceExporter, err := initTracesExporter(context.Background())
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithResource(resource),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSpanProcessor(
			tracesdk.NewBatchSpanProcessor(
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

	defer tp.Shutdown(context.Background())

	// metrics
	metricsExporter, err := initMetricsExporter(context.Background())
	if err != nil {
		return err
	}

	mp := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(
			metricsdk.NewPeriodicReader(
				metricsExporter,
				metricsdk.WithInterval(15*time.Second),
				metricsdk.WithProducer(runtime.NewProducer()),
			),
		),
	)

	otel.SetMeterProvider(mp)

	defer mp.Shutdown(context.Background())

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
		slog.InfoContext(context.Background(), fmt.Sprintf("http server listening on %s", config.HttpAddress()))
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
		slog.ErrorContext(context.Background(), "failed to start server", "error", err)
		return err
	}

	// graceful shutdown
	wait := make(chan struct{})

	go func() {
		defer close(wait)
		wg.Wait()
	}()

	close(cacheStop)

	slog.InfoContext(context.Background(), "successfully stopped cache")

	close(exportStop)

	slog.InfoContext(context.Background(), "successfully stopped export")

	select {
	case <-wait:
	case <-time.After(30 * time.Second):
	}

	slog.InfoContext(context.Background(), "successfully stopped server")

	return nil
}

func initLogsExporter(ctx context.Context) (logsdk.Exporter, error) {
	switch config.LogsExporter() {
	case "dd-otlp":
		return otlploghttp.New(
			ctx,
			otlploghttp.WithEndpointURL(config.LogsAddress()),
			otlploghttp.WithURLPath(config.LogsUrlPath()),
			otlploghttp.WithHeaders(
				map[string]string{
					"dd-protocol": "otlp",
					"dd-api-key":  config.LogsAPIToken(),
				},
			),
		)
	default:
		return stdoutlog.New()
	}
}

func initTracesExporter(ctx context.Context) (tracesdk.SpanExporter, error) {
	switch config.TracesExporter() {
	default:
		return otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(config.TracesAddress()),
			otlptracehttp.WithInsecure(),
		)
	}
}

func initMetricsExporter(ctx context.Context) (metricsdk.Exporter, error) {
	switch config.MetricsExporter() {
	default:
		return otlpmetrichttp.New(
			ctx,
			otlpmetrichttp.WithEndpoint(config.MetricsAddress()),
			otlpmetrichttp.WithInsecure(),
		)
	}
}

func initWriteClient() writer.Writer {
	switch config.WriteClient() {
	case "postgres":
		return postgreswriter.NewWriter(
			writer.WithLocation(config.WriteClientLocation()),
		)
	case "dynamodb":
		return dynamodbwriter.NewWriter(
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
	case "dynamodb":
		return dynamodbreader.NewReader(
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
