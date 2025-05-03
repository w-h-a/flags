package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Parser struct {
}

func (p *Parser) ParseFlagKey(ctx context.Context, r *http.Request) (string, error) {
	vars := mux.Vars(r)

	flagKey := vars["key"]

	if len(flagKey) == 0 {
		return "", fmt.Errorf("flag key is required")
	}

	return flagKey, nil
}
