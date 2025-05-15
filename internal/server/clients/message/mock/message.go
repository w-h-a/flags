package mock

import (
	"context"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/message"
)

type Client struct {
	options   message.Options
	wasCalled bool
	mtx       sync.RWMutex
}

func (c *Client) Send(ctx context.Context, diff message.Diff) error {
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

func NewMessageClient(opts ...message.Option) message.Client {
	options := message.NewOptions(opts...)

	c := &Client{
		options: options,
		mtx:     sync.RWMutex{},
	}

	return c
}
