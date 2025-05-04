package local

import (
	"context"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/message"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options message.Options
}

func (c *client) Send(ctx context.Context, diff message.Diff, wg *sync.WaitGroup) error {
	defer wg.Done()

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

func NewMessageClient(opts ...message.Option) message.Client {
	options := message.NewOptions(opts...)

	c := &client{
		options: options,
	}

	return c
}
