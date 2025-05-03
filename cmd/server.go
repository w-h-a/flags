package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/flags/internal/server"
	"github.com/w-h-a/flags/internal/server/clients/file"
	localfile "github.com/w-h-a/flags/internal/server/clients/file/local"
	localmessage "github.com/w-h-a/flags/internal/server/clients/message/local"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/notify"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
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

	// clients
	// TODO: get this from config
	fileClient := localfile.NewFileClient(
		// TODO: get this from config
		file.WithDir("."),
		file.WithFiles("/flags.yaml"),
	)

	messageClient := localmessage.NewMessageClient()

	// server
	httpServer, cacheService, notifyService, err := server.Factory(fileClient, messageClient)
	if err != nil {
		return err
	}

	// wait group and error chan
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 2)

	// start http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- httpServer.Start()
	}()

	// start cache updater
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- updateCache(cacheService, notifyService, stop)
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

	select {
	case <-wait:
	case <-time.After(30 * time.Second):
	}

	log.Info("successfully stopped server")

	close(stop)

	log.Info("successfully stopped cache")

	return nil
}

func updateCache(cacheService *cache.Service, notifyService *notify.Service, stop chan struct{}) error {
	// TODO: confirm poll interval is valid

	// TODO: retrieve from config
	ticker := time.NewTicker(time.Minute)

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
