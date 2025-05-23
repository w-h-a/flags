package local

import (
	"context"
	"os"

	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/pkg/telemetry/log"
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
		log.Fatalf("failed to configure local file client: %v", err)
	}

	c := &client{
		options: options,
	}

	return c
}
