package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w-h-a/flags/internal/server/services/cache"
)

type Flags struct {
	cacheService *cache.Service
}

func (f *Flags) GetAll(w http.ResponseWriter, r *http.Request) {
	flags := f.cacheService.Flags()

	bs, err := json.Marshal(flags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func NewFlagsHandler(cacheService *cache.Service) *Flags {
	return &Flags{
		cacheService: cacheService,
	}
}
