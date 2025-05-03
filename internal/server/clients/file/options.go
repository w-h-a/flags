package file

import "context"

type Option func(o *Options)

type Options struct {
	Dir     string
	Files   []string
	Context context.Context
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
