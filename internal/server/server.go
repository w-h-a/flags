package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/w-h-a/flags/internal/server/clients/exporter"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/config"
	httphandlers "github.com/w-h-a/flags/internal/server/handlers/http"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/export"
	"github.com/w-h-a/flags/internal/server/services/notify"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"github.com/w-h-a/pkg/telemetry/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func Factory(
	readClient reader.Reader,
	exportClient exporter.Exporter,
	notifyClient notifier.Notifier,
) (serverv2.Server, *cache.Service, *export.Service, *notify.Service, error) {
	// services
	cacheService := cache.New(readClient)
	exportService := export.New(exportClient)
	notifyService := notify.New(notifyClient)

	old, new, err := cacheService.RetrieveFlags()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	notifyService.Notify(old, new)

	// base server options
	opts := []serverv2.ServerOption{
		serverv2.ServerWithNamespace(config.Env()),
		serverv2.ServerWithName(config.Name()),
		serverv2.ServerWithVersion(config.Version()),
	}

	// create http server
	router := mux.NewRouter()

	httpOFREP := httphandlers.NewOFREPHandler(cacheService, exportService)

	router.Methods(http.MethodPost).Path("/ofrep/v1/evaluate/flags/{key}").HandlerFunc(httpOFREP.PostOne)
	router.Methods(http.MethodPost).Path("/ofrep/v1/evaluate/flags").HandlerFunc(httpOFREP.PostAll)
	router.Methods(http.MethodGet).Path("/ofrep/v1/configuration").HandlerFunc(httpOFREP.GetConfig)

	httpStatus := httphandlers.NewStatusHandler(cacheService)

	router.Methods(http.MethodGet).Path("/status").HandlerFunc(httpStatus.GetStatus)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
		httpserver.HttpServerWithMiddleware(httphandlers.NewAuthMiddleware()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	handler := otelhttp.NewHandler(
		router,
		"",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string { return r.URL.Path }),
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
		otelhttp.WithFilter(func(r *http.Request) bool { return r.URL.Path != "/status" }),
	)

	httpServer.Handle(handler)

	return httpServer, cacheService, exportService, notifyService, nil
}

func UpdateCache(
	cacheService *cache.Service,
	notifyService *notify.Service,
	stop chan struct{},
	dur time.Duration,
) error {
	ticker := time.NewTicker(dur)

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

func ExportReports(
	exportService *export.Service,
	stop chan struct{},
	dur time.Duration,
) error {
	ticker := time.NewTicker(dur)

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
