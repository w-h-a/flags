package reader

import (
	"context"
	"fmt"
)

type Option func(o *Options)

type Options struct {
	Dir     string
	File    string
	Token   string
	Context context.Context
}

func (o Options) Validate() error {
	if len(o.Dir) == 0 {
		return fmt.Errorf("missing directory")
	}

	if len(o.File) == 0 {
		return fmt.Errorf("missing file")
	}

	return nil
}

func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

func WithFile(file string) Option {
	return func(o *Options) {
		o.File = file
	}
}

func WithToken(token string) Option {
	return func(o *Options) {
		o.Token = token
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
