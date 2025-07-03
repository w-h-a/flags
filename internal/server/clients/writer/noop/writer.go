package noop

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
		slog.ErrorContext(context.Background(), fmt.Sprintf("failed to validate noop writer: %v", err))
		os.Exit(1)
	}

	c := &client{
		options: options,
	}

	return c
}
