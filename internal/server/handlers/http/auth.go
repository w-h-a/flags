package http

import (
	"net/http"
	"strings"

	"github.com/w-h-a/flags/internal/server/config"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
)

const (
	BearerScheme = "Bearer "
)

type AuthMiddleware struct {
	handler http.Handler
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/status" {
		m.handler.ServeHTTP(w, r)
		return
	}

	errBody := map[string]any{"error": "not authenticated"}

	token := ""

	authHeader := r.Header.Get("authorization")
	if len(authHeader) > 0 {
		if !strings.HasPrefix(authHeader, BearerScheme) {
			writeRsp(w, http.StatusUnauthorized, errBody)
			return
		}
		token = authHeader[len(BearerScheme):]
	}

	if !config.CheckAPIKey(token) {
		writeRsp(w, http.StatusUnauthorized, errBody)
		return
	}

	m.handler.ServeHTTP(w, r)
}

func NewAuthMiddleware() httpserver.Middleware {
	return func(h http.Handler) http.Handler {
		return &AuthMiddleware{h}
	}
}
