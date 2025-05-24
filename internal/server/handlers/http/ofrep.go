package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/export"
)

type OFREP struct {
	cacheService  *cache.Service
	exportService *export.Service
	parser        *Parser
}

func (o *OFREP) PostOne(w http.ResponseWriter, r *http.Request) {
	ctx := RequestToContext(r)

	flagKey, err := o.parser.ParseFlagKey(ctx, r)
	if err != nil {
		bs, _ := json.Marshal(map[string]any{"error": err.Error()})
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(bs))
		return
	}

	flagState, err := o.cacheService.EvaluateFlag(ctx, flagKey)
	if err != nil && errors.Is(err, cache.ErrNotFound) {
		bs, _ := json.Marshal(flagState)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(bs))
		return
	} else if err != nil {
		bs, _ := json.Marshal(map[string]any{"error": err.Error()})
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(bs))
		return
	}

	if config.ExportReports() {
		event := export.Event{
			CreationDate: time.Now().Unix(),
			Key:          flagState.Key,
			Value:        flagState.Value,
			Variant:      flagState.Variant,
			Reason:       flagState.Reason,
			ErrorCode:    flagState.ErrorCode,
			ErrorMessage: flagState.ErrorMessage,
		}

		o.exportService.Add(event)
	}

	bs, _ := json.Marshal(flagState)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func (o *OFREP) PostAll(w http.ResponseWriter, r *http.Request) {
	ctx := RequestToContext(r)

	flags := o.cacheService.EvaluateFlags(ctx)

	bs, _ := json.Marshal(flags)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func (o *OFREP) GetConfig(w http.ResponseWriter, r *http.Request) {
	rsp := map[string]any{
		"name": config.Name(),
		"capabilities": map[string]any{
			"cacheInvalidation": map[string]any{
				"polling": map[string]any{
					"enabled": true,
				},
			},
			"flagEvaluation": map[string]any{
				"supportedTypes": []string{"int", "float", "string", "boolean"},
			},
		},
	}

	bs, _ := json.Marshal(rsp)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func NewOFREPHandler(
	cacheService *cache.Service,
	exportService *export.Service,
) *OFREP {
	return &OFREP{
		cacheService:  cacheService,
		exportService: exportService,
		parser:        &Parser{},
	}
}
