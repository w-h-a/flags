package mock

import (
	"context"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/file"
)

type initialFlagsKey struct{}

func WithInitialFlags(flags map[string]*flags.Flag) file.Option {
	return func(o *file.Options) {
		o.Context = context.WithValue(o.Context, initialFlagsKey{}, flags)
	}
}

func InitialFlags(ctx context.Context) (map[string]*flags.Flag, bool) {
	fs, ok := ctx.Value(initialFlagsKey{}).(map[string]*flags.Flag)
	return fs, ok
}

type updatedFlagsKey struct{}

func WithUpdatedFlags(flags map[string]*flags.Flag) file.Option {
	return func(o *file.Options) {
		o.Context = context.WithValue(o.Context, updatedFlagsKey{}, flags)
	}
}

func UpdatedFlags(ctx context.Context) (map[string]*flags.Flag, bool) {
	fs, ok := ctx.Value(updatedFlagsKey{}).(map[string]*flags.Flag)
	return fs, ok
}
