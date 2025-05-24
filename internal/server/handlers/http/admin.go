package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/w-h-a/flags/internal/server/services/admin"
)

type Admin struct {
	adminService *admin.Service
	parser       *Parser
}

func (a *Admin) GetOne(w http.ResponseWriter, r *http.Request) {

}

func (a *Admin) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (a *Admin) PutOne(w http.ResponseWriter, r *http.Request) {
	ctx := RequestToContext(r)

	flag, err := a.parser.ParsePutOneBody(ctx, r)
	if err != nil {
		bs, _ := json.Marshal(map[string]any{"error": err.Error()})
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(bs))
		return
	}

	if len(flag) != 1 {
		bs, _ := json.Marshal(map[string]any{"error": "body does not contain exactly one flag"})
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(bs))
		return
	}

	flagKey := ""

	for k := range flag {
		flagKey = k
	}

	found := true

	if _, err := a.adminService.RetrieveFlag(ctx, flagKey); err != nil && errors.Is(err, admin.ErrNotFound) {
		found = false
	} else if err != nil {
		bs, _ := json.Marshal(map[string]any{"error": err.Error()})
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(bs))
		return
	}

	upserted, err := a.adminService.UpsertFlag(ctx, flagKey, flag)
	if err != nil {
		bs, _ := json.Marshal(map[string]any{"error": err.Error()})
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(bs))
		return
	}

	bs, _ := json.Marshal(upserted)
	w.Header().Set("content-type", "application/json")
	if found {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	fmt.Fprint(w, string(bs))
}

func (a *Admin) PatchOne(w http.ResponseWriter, r *http.Request) {

}

func NewAdminHandler(adminService *admin.Service) *Admin {
	return &Admin{
		adminService: adminService,
		parser:       &Parser{},
	}
}
