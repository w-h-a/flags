package http

import (
	"errors"
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
	ctx := reqToCtx(r)

	flagKey, err := o.parser.ParseFlagKey(ctx, r)
	if err != nil {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	evalCtx, err := o.parser.ParsePostOneBody(ctx, r)
	if err != nil {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	flagState, err := o.cacheService.EvaluateFlag(ctx, flagKey, evalCtx)
	if err != nil && errors.Is(err, cache.ErrNotFound) {
		writeRsp(w, http.StatusNotFound, flagState)
		return
	} else if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
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

	writeRsp(w, http.StatusOK, flagState)
}

func (o *OFREP) PostAll(w http.ResponseWriter, r *http.Request) {
	ctx := reqToCtx(r)

	flags := o.cacheService.EvaluateFlags(ctx)

	writeRsp(w, http.StatusOK, flags)
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
