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
	flag, err := a.parser.ParsePutOneBody(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	if len(flag) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: body does not contain exactly one flag")
		return
	}

	flagKey := ""

	for k := range flag {
		flagKey = k
	}

	found := true

	if _, err := a.adminService.RetrieveFlag(r.Context(), flagKey); err != nil && errors.Is(err, admin.ErrNotFound) {
		found = false
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	upserted, err := a.adminService.UpsertFlag(r.Context(), flagKey, flag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	bs, err := json.Marshal(upserted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

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
