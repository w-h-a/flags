package exporter

import (
	"context"
	"fmt"
)

type Option func(o *Options)

type Options struct {
	Dir     string
	Format  string
	Context context.Context
}

func (o Options) Validate() error {
	if len(o.Dir) == 0 {
		return fmt.Errorf("missing directory")
	}

	return nil
}

func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

func WithFormat(format string) Option {
	return func(o *Options) {
		o.Format = format
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Format:  "csv",
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
