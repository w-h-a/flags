package local

import (
	"context"
	"log/slog"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
)

type client struct {
	options notifier.Options
}

func (c *client) Notify(ctx context.Context, diff flags.Diff) error {
	for k := range diff.Deleted {
		slog.InfoContext(ctx, "flag removed", "flag", k)
	}

	for k := range diff.Added {
		slog.InfoContext(ctx, "flag added", "flag", k)
	}

	for k, v := range diff.Updated {
		if v.After.IsDisabled() != v.Before.IsDisabled() {
			if v.After.IsDisabled() {
				slog.InfoContext(ctx, "flag is OFF", "flag", k)
				continue
			}
			slog.InfoContext(ctx, "flag is ON", "flag", k)
			continue
		}
		slog.InfoContext(ctx, "flag is updated", "flag", k)
	}

	return nil
}

func NewNotifier(opts ...notifier.Option) notifier.Notifier {
	options := notifier.NewOptions(opts...)

	c := &client{
		options: options,
	}

	return c
}
