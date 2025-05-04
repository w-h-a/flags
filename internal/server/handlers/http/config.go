package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w-h-a/flags/internal/server/config"
)

type Config struct {
}

func (c *Config) GetConfig(w http.ResponseWriter, r *http.Request) {
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

	bs, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

func NewConfigHandler() *Config {
	return &Config{}
}
