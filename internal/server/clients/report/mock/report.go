package mock

import (
	"context"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/report"
	"github.com/w-h-a/pkg/telemetry/log"
)

type Client struct {
	options report.Options
	parser  *report.Parser
	records []report.Record
	mtx     sync.RWMutex
}

func (c *Client) Create(ctx context.Context, records []report.Record) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.records = append(c.records, records...)

	return nil
}

func (c *Client) Records() []report.Record {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.records
}

func NewReportClient(opts ...report.Option) report.Client {
	options := report.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatal(err)
	}

	p := &report.Parser{}

	p.FilenameTemplate = p.ParseTemplate("filenameTemplate", report.FilenameTemplate)
	p.CsvTemplate = p.ParseTemplate("csvTemplate", report.CsvTemplate)

	c := &Client{
		options: options,
		parser:  p,
		mtx:     sync.RWMutex{},
	}

	return c
}
