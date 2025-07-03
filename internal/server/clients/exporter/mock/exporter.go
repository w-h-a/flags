package mock

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/exporter"
)

type Client struct {
	options exporter.Options
	parser  *exporter.Parser
	records []exporter.Record
	mtx     sync.RWMutex
}

func (c *Client) Export(ctx context.Context, records []exporter.Record) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.records = append(c.records, records...)

	return nil
}

func (c *Client) Records() []exporter.Record {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.records
}

func NewExporter(opts ...exporter.Option) exporter.Exporter {
	options := exporter.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		slog.ErrorContext(context.Background(), "failed to validate mock exporter options", "error", err)
		os.Exit(1)
	}

	p := &exporter.Parser{}

	p.FilenameTemplate = p.ParseTemplate("filenameTemplate", exporter.FilenameTemplate)
	p.CsvTemplate = p.ParseTemplate("csvTemplate", exporter.CsvTemplate)

	c := &Client{
		options: options,
		parser:  p,
		mtx:     sync.RWMutex{},
	}

	return c
}
