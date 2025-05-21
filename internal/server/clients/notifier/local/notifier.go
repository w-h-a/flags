package local

import (
	"context"

	"github.com/w-h-a/flags/internal/server/clients/notifier"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options notifier.Options
}

func (c *client) Notify(ctx context.Context, diff notifier.Diff) error {
	for k := range diff.Deleted {
		log.Infof("flag %v removed", k)
	}

	for k := range diff.Added {
		log.Infof("flag %v added", k)
	}

	for k, v := range diff.Updated {
		if v.After.IsDisabled() != v.Before.IsDisabled() {
			if v.After.IsDisabled() {
				log.Infof("flag %v is OFF", k)
				continue
			}
			log.Infof("flag %v is ON", k)
			continue
		}
		log.Infof("flag %v updated", k)
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
