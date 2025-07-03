package local

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
)

type client struct {
	options notifier.Options
}

func (c *client) Notify(ctx context.Context, diff flags.Diff) error {
	for k := range diff.Deleted {
		slog.InfoContext(ctx, fmt.Sprintf("flag %v removed", k))
	}

	for k := range diff.Added {
		slog.InfoContext(ctx, fmt.Sprintf("flag %v added", k))
	}

	for k, v := range diff.Updated {
		if v.After.IsDisabled() != v.Before.IsDisabled() {
			if v.After.IsDisabled() {
				slog.InfoContext(ctx, fmt.Sprintf("flag %v is OFF", k))
				continue
			}
			slog.InfoContext(ctx, fmt.Sprintf("flag %v is ON", k))
			continue
		}
		slog.InfoContext(ctx, fmt.Sprintf("flag %v is updated", k))
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
