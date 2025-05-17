package local

import (
	"context"
	"os"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options file.Options
}

func (c *client) Read(ctx context.Context) (map[string]*flags.Flag, error) {
	// TODO: generalize
	file := c.options.Files[0]

	bs, err := os.ReadFile(c.options.Dir + file)
	if err != nil {
		return nil, err
	}

	return flags.Factory(bs, config.FlagFormat())
}

func NewFileClient(opts ...file.Option) file.Client {
	options := file.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatalf("failed to configure local file client: %v", err)
	}

	c := &client{
		options: options,
	}

	return c
}
