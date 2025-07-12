package local

import (
	"context"
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
		detail := "failed to configure local file reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	c := &client{
		options: options,
	}

	return c
}
