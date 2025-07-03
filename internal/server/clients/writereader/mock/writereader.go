package mock

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/clients/writereader"
)

type client struct {
	options writereader.Options
	store   map[string][]byte
	mtx     sync.RWMutex
}

func (c *client) Write(ctx context.Context, key string, bs []byte) error {
	if err, ok := ctx.Value("error_write").(string); ok {
		return fmt.Errorf("%s", err)
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.store[key] = bs

	return nil
}

func (c *client) ReadByKey(ctx context.Context, key string) ([]byte, error) {
	if err, ok := ctx.Value("error_read_by_key").(string); ok {
		return nil, fmt.Errorf("%s", err)
	}

	c.mtx.RLock()
	defer c.mtx.RUnlock()

	bs, found := c.store[key]
	if !found {
		return nil, reader.ErrRecordNotFound
	}

	return bs, nil
}

func (c *client) Read(ctx context.Context) ([]byte, error) {
	if err, ok := ctx.Value("error_read").(string); ok {
		return nil, fmt.Errorf("%s", err)
	}

	c.mtx.RLock()
	defer c.mtx.RUnlock()

	bs := []byte{}

	for _, v := range c.store {
		bs = append(bs, []byte("\n")...)
		bs = append(bs, v...)
	}

	return bs, nil
}

func NewWriteReader(opts ...writereader.Option) writereader.WriteReader {
	options := writereader.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		slog.ErrorContext(context.Background(), "failed to validate mock write reader options", "error", err)
		os.Exit(1)
	}

	c := &client{
		options: options,
		store:   map[string][]byte{},
		mtx:     sync.RWMutex{},
	}

	return c
}
