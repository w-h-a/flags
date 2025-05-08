package file

import (
	"context"
	"fmt"
)

type Option func(o *Options)

type Options struct {
	Dir     string
	Files   []string
	Token   string
	Context context.Context
}

func (o Options) Validate() error {
	if len(o.Dir) == 0 {
		return fmt.Errorf("missing directory")
	}

	if len(o.Files) == 0 {
		return fmt.Errorf("missing files")
	}

	return nil
}

func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

func WithFiles(files ...string) Option {
	return func(o *Options) {
		o.Files = append(o.Files, files...)
	}
}

func WithToken(token string) Option {
	return func(o *Options) {
		o.Token = token
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Files:   []string{},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
