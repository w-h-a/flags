package local

import (
	"context"
	"os"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options reader.Options
}

func (c *client) Read(ctx context.Context) (map[string]*flags.Flag, error) {
	bs, err := os.ReadFile(c.options.Location)
	if err != nil {
		return nil, err
	}

	return flags.Factory(bs, config.FlagFormat())
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
