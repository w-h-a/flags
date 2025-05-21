package reader

import (
	"context"

	"github.com/w-h-a/flags/internal/flags"
)

type Reader interface {
	Read(ctx context.Context) (map[string]*flags.Flag, error)
}
