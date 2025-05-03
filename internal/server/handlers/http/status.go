package http

import (
	"encoding/json"
	"fmt"
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

	bs, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func NewStatusHandler(cacheService *cache.Service) *Status {
	return &Status{
		cacheService: cacheService,
	}
}
