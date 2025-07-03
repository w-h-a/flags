package local

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/w-h-a/flags/internal/server/clients/reader"
)

type client struct {
	options reader.Options
}

func (c *client) ReadByKey(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (c *client) Read(ctx context.Context) ([]byte, error) {
	return os.ReadFile(c.options.Location)
}

func NewReader(opts ...reader.Option) reader.Reader {
	options := reader.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		slog.ErrorContext(context.Background(), fmt.Sprintf("failed to configure local file reader: %v", err))
		os.Exit(1)
	}

	c := &client{
		options: options,
	}

	return c
}
