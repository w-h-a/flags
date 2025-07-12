package noop

import (
	"context"
	"log/slog"

	"github.com/w-h-a/flags/internal/server/clients/writer"
)

type client struct {
	options writer.Options
}

func (c *client) Write(ctx context.Context, key string, bs []byte) error {
	return nil
}

func NewWriter(opts ...writer.Option) writer.Writer {
	options := writer.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		detail := "failed to validate noop writer"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	c := &client{
		options: options,
	}

	return c
}
