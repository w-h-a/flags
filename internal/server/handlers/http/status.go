package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w-h-a/flags/internal/server/services/monitor"
)

type Status struct {
	monitorService *monitor.Service
}

func (s *Status) GetStatus(w http.ResponseWriter, r *http.Request) {
	bs, err := json.Marshal(s.monitorService.Info())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func NewStatusHandler(monitorService *monitor.Service) *Status {
	return &Status{
		monitorService: monitorService,
	}
}
