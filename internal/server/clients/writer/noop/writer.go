package noop

import (
	"context"

	"github.com/w-h-a/flags/internal/server/clients/writer"
	"github.com/w-h-a/pkg/telemetry/log"
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
		log.Fatal(err)
	}

	c := &client{
		options: options,
	}

	return c
}
