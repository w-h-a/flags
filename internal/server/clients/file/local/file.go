package local

import (
	"context"
	"os"

	"github.com/w-h-a/flags/internal/server/clients/file"
)

type client struct {
	options file.Options
	parser  *file.Parser
}

func (c *client) Read(ctx context.Context) (map[string]*file.Flag, error) {
	file := "/"

	// TODO: generalize
	if len(c.options.Files) > 0 {
		file = c.options.Files[0]
	}

	bs, err := os.ReadFile(c.options.Dir + file)
	if err != nil {
		return nil, err
	}

	return c.parser.ParseFlags(bs)
}

func NewFileClient(opts ...file.Option) file.Client {
	options := file.NewOptions(opts...)

	c := &client{
		options: options,
		parser:  &file.Parser{},
	}

	return c
}
