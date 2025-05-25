package http

import (
	"net/http"

	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/flags/internal/server/services/cache"
)

type Status struct {
	cacheService *cache.Service
}

func (s *Status) GetStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]any{
		"env":        config.Env(),
		"name":       config.Name(),
		"version":    config.Version(),
		"lastUpdate": s.cacheService.LastUpdate(),
	}

	writeRsp(w, http.StatusOK, status)
}

func NewStatusHandler(cacheService *cache.Service) *Status {
	return &Status{
		cacheService: cacheService,
	}
}
