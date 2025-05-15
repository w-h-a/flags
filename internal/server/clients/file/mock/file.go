package mock

import (
	"context"
	"maps"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/pkg/telemetry/log"
)

type Client struct {
	options      file.Options
	initialFlags map[string]*file.Flag
	updatedFlags map[string]*file.Flag
	callCount    int
	mtx          sync.RWMutex
}

func (c *Client) Read(ctx context.Context) (map[string]*file.Flag, error) {
	c.mtx.RLock()

	callCount := c.callCount

	result := map[string]*file.Flag{}

	if callCount == 0 {
		maps.Copy(result, c.initialFlags)
	} else {
		maps.Copy(result, c.updatedFlags)
	}

	c.mtx.RUnlock()

	c.mtx.Lock()
	c.callCount++
	c.mtx.Unlock()

	return result, nil
}

func (c *Client) CallCount() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.callCount
}

func NewFileClient(opts ...file.Option) file.Client {
	options := file.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatal(err)
	}

	c := &Client{
		options: options,
		mtx:     sync.RWMutex{},
	}

	if fs, ok := InitialFlags(options.Context); ok {
		c.initialFlags = fs
	}

	if fs, ok := UpdatedFlags(options.Context); ok {
		c.updatedFlags = fs
	}

	return c
}
