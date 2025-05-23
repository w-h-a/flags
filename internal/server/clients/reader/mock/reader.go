package mock

import (
	"context"
	"maps"
	"sync"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/pkg/telemetry/log"
	"gopkg.in/yaml.v3"
)

type Client struct {
	options      reader.Options
	initialFlags map[string]*flags.Flag
	updatedFlags map[string]*flags.Flag
	callCount    int
	mtx          sync.RWMutex
}

func (c *Client) ReadByKey(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (c *Client) Read(ctx context.Context) ([]byte, error) {
	c.mtx.RLock()

	callCount := c.callCount

	result := map[string]*flags.Flag{}

	if callCount == 0 {
		maps.Copy(result, c.initialFlags)
	} else {
		maps.Copy(result, c.updatedFlags)
	}

	c.mtx.RUnlock()

	c.mtx.Lock()
	c.callCount++
	c.mtx.Unlock()

	return yaml.Marshal(result)
}

func (c *Client) CallCount() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.callCount
}

func NewReader(opts ...reader.Option) reader.Reader {
	options := reader.NewOptions(opts...)

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
