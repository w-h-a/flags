package http

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/flags/internal/flags"
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

func (p *Parser) ParsePutOneBody(ctx context.Context, r *http.Request) (map[string]*flags.Flag, error) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return flags.Factory(bs, "json")
}
