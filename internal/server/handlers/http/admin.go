package http

import (
	"errors"
	"net/http"
	"sync"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/services/admin"
)

type Admin struct {
	adminService *admin.Service
	parser       *Parser
	mtx          sync.RWMutex
}

func (a *Admin) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx := reqToCtx(r)

	flagKey, err := a.parser.ParseFlagKey(ctx, r)
	if err != nil {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	flag, err := a.adminService.RetrieveFlag(ctx, flagKey)
	if err != nil && errors.Is(err, flags.ErrNotFound) {
		writeRsp(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	} else if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeRsp(w, http.StatusOK, flag)
}

func (a *Admin) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := reqToCtx(r)

	flags, err := a.adminService.RetrieveFlags(ctx)
	if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeRsp(w, http.StatusOK, flags)
}

func (a *Admin) PutOne(w http.ResponseWriter, r *http.Request) {
	ctx := reqToCtx(r)

	flag, err := a.parser.ParsePutOneBody(ctx, r)
	if err != nil {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	if len(flag) != 1 {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": "body does not contain exactly one flag"})
		return
	}

	flagKey := ""

	for k := range flag {
		flagKey = k
	}

	found := true

	if _, err := a.adminService.RetrieveFlag(ctx, flagKey); err != nil && errors.Is(err, flags.ErrNotFound) {
		found = false
	} else if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	upserted, err := a.adminService.UpsertFlag(ctx, flagKey, flag)
	if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	if found {
		writeRsp(w, http.StatusOK, upserted)
	} else {
		writeRsp(w, http.StatusCreated, upserted)
	}
}

func (a *Admin) PatchOne(w http.ResponseWriter, r *http.Request) {
	ctx := reqToCtx(r)

	flagKey, err := a.parser.ParseFlagKey(ctx, r)
	if err != nil {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	disabledPatch, err := a.parser.ParsePatchOneBody(ctx, r)
	if err != nil {
		writeRsp(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	flag, err := a.adminService.RetrieveFlag(ctx, flagKey)
	if err != nil && errors.Is(err, flags.ErrNotFound) {
		writeRsp(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	} else if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	a.mtx.Lock()
	f := flag[flagKey]
	f.Disabled = disabledPatch.Disabled
	a.mtx.Unlock()

	upserted, err := a.adminService.UpsertFlag(ctx, flagKey, flag)
	if err != nil {
		writeRsp(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeRsp(w, http.StatusOK, upserted)
}

func NewAdminHandler(adminService *admin.Service) *Admin {
	return &Admin{
		adminService: adminService,
		parser:       &Parser{},
		mtx:          sync.RWMutex{},
	}
}
