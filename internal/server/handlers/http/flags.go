package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/w-h-a/flags/internal/server/services/cache"
)

type Flags struct {
	cacheService *cache.Service
	parser       *Parser
}

func (f *Flags) PostOne(w http.ResponseWriter, r *http.Request) {
	flagKey, err := f.parser.ParseFlagKey(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	flagValue, err := f.cacheService.Flag(flagKey)
	if err != nil && errors.Is(err, cache.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "error: %s %v", flagKey, err)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	bs, err := json.Marshal(flagValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func (f *Flags) PostAll(w http.ResponseWriter, r *http.Request) {
	flags := f.cacheService.Flags()

	bs, err := json.Marshal(flags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func NewFlagsHandler(cacheService *cache.Service) *Flags {
	return &Flags{
		cacheService: cacheService,
		parser:       &Parser{},
	}
}
