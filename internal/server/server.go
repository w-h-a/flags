package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/flags/internal/server/clients/message"
	"github.com/w-h-a/flags/internal/server/clients/report"
	"github.com/w-h-a/flags/internal/server/config"
	httphandlers "github.com/w-h-a/flags/internal/server/handlers/http"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/export"
	"github.com/w-h-a/flags/internal/server/services/notify"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func Factory(
	fileClient file.Client,
	reportClient report.Client,
	notifiers ...message.Client,
) (serverv2.Server, *cache.Service, *export.Service, *notify.Service, error) {
	// services
	cacheService := cache.New(fileClient)
	exportService := export.New(reportClient)
	notifyService := notify.New(notifiers...)

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

	httpFlags := httphandlers.NewFlagsHandler(cacheService, exportService)
	httpConfig := httphandlers.NewConfigHandler()
	httpStatus := httphandlers.NewStatusHandler(cacheService)

	router.Methods(http.MethodPost).Path("/ofrep/v1/evaluate/flags/{key}").HandlerFunc(httpFlags.PostOne)
	router.Methods(http.MethodPost).Path("/ofrep/v1/evaluate/flags").HandlerFunc(httpFlags.PostAll)
	router.Methods(http.MethodGet).Path("/ofrep/v1/configuration").HandlerFunc(httpConfig.GetConfig)
	router.Methods(http.MethodGet).Path("/status").HandlerFunc(httpStatus.GetStatus)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
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
