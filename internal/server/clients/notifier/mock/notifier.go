package mock

import (
	"context"
	"sync"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
)

type Client struct {
	options   notifier.Options
	wasCalled bool
	mtx       sync.RWMutex
}

func (c *Client) Notify(ctx context.Context, diff flags.Diff) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.wasCalled = true

	return nil
}

func (c *Client) WasCalled() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.wasCalled
}

func NewNotifier(opts ...notifier.Option) notifier.Notifier {
	options := notifier.NewOptions(opts...)

	c := &Client{
		options: options,
		mtx:     sync.RWMutex{},
	}

	return c
}
