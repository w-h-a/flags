package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/export"
)

type Flags struct {
	cacheService  *cache.Service
	exportService *export.Service
	parser        *Parser
}

func (f *Flags) PostOne(w http.ResponseWriter, r *http.Request) {
	flagKey, err := f.parser.ParseFlagKey(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	flagState, err := f.cacheService.Flag(flagKey)
	if err != nil && errors.Is(err, cache.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "error: %s %v", flagKey, err)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	event := export.Event{
		CreationDate: time.Now().Unix(),
		Key:          flagState.Key,
		Value:        flagState.Value,
		Variant:      flagState.Variant,
		Reason:       flagState.Reason,
		ErrorCode:    flagState.ErrorCode,
	}

	f.exportService.Add(event)

	bs, err := json.Marshal(flagState)
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

func NewFlagsHandler(cacheService *cache.Service, exportService *export.Service) *Flags {
	return &Flags{
		cacheService:  cacheService,
		exportService: exportService,
		parser:        &Parser{},
	}
}
