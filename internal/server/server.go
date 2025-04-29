package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/flags/internal/server/config"
	httphandlers "github.com/w-h-a/flags/internal/server/handlers/http"
	"github.com/w-h-a/flags/internal/server/services/monitor"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func Factory() serverv2.Server {
	// services
	monitorService := monitor.New()

	// base server options
	opts := []serverv2.ServerOption{
		serverv2.ServerWithNamespace(config.Env()),
		serverv2.ServerWithName(config.Name()),
		serverv2.ServerWithVersion(config.Version()),
	}

	// create http server
	router := mux.NewRouter()

	httpStatus := httphandlers.NewStatusHandler(monitorService)

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

	return httpServer
}
