package file

import (
	"context"

	"github.com/w-h-a/flags/internal/flags"
)

type Client interface {
	Read(ctx context.Context) (map[string]*flags.Flag, error)
}
