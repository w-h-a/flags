package mock

import (
	"context"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/reader"
)

type initialFlagsKey struct{}

func WithInitialFlags(flags map[string]*flags.Flag) reader.Option {
	return func(o *reader.Options) {
		o.Context = context.WithValue(o.Context, initialFlagsKey{}, flags)
	}
}

func InitialFlags(ctx context.Context) (map[string]*flags.Flag, bool) {
	fs, ok := ctx.Value(initialFlagsKey{}).(map[string]*flags.Flag)
	return fs, ok
}

type updatedFlagsKey struct{}

func WithUpdatedFlags(flags map[string]*flags.Flag) reader.Option {
	return func(o *reader.Options) {
		o.Context = context.WithValue(o.Context, updatedFlagsKey{}, flags)
	}
}

func UpdatedFlags(ctx context.Context) (map[string]*flags.Flag, bool) {
	fs, ok := ctx.Value(updatedFlagsKey{}).(map[string]*flags.Flag)
	return fs, ok
}
