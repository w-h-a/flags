package mock

import (
	"context"

	"github.com/w-h-a/flags/internal/server/clients/file"
)

type initialFlagsKey struct{}

func WithInitialFlags(flags map[string]*file.Flag) file.Option {
	return func(o *file.Options) {
		o.Context = context.WithValue(o.Context, initialFlagsKey{}, flags)
	}
}

func InitialFlags(ctx context.Context) (map[string]*file.Flag, bool) {
	fs, ok := ctx.Value(initialFlagsKey{}).(map[string]*file.Flag)
	return fs, ok
}

type updatedFlagsKey struct{}

func WithUpdatedFlags(flags map[string]*file.Flag) file.Option {
	return func(o *file.Options) {
		o.Context = context.WithValue(o.Context, updatedFlagsKey{}, flags)
	}
}

func UpdatedFlags(ctx context.Context) (map[string]*file.Flag, bool) {
	fs, ok := ctx.Value(updatedFlagsKey{}).(map[string]*file.Flag)
	return fs, ok
}
