package writer

import (
	"context"
	"fmt"
)

type Option func(o *Options)

type Options struct {
	Location string
	Token    string
	Context  context.Context
}

func (o Options) Validate() error {
	if len(o.Location) == 0 {
		return fmt.Errorf("missing flag location")
	}

	return nil
}

func WithLocation(location string) Option {
	return func(o *Options) {
		o.Location = location
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
