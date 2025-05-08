package local

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/w-h-a/flags/internal/server/clients/report"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options report.Options
	parser  *report.Parser
}

func (c *client) Create(ctx context.Context, records []report.Record) error {
	filename, err := c.parser.ParseFilename(c.options.Format)
	if err != nil {
		return err
	}

	filePath := c.options.Dir + "/" + filename

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, record := range records {
		var line []byte
		var err error

		switch strings.ToLower(c.options.Format) {
		case "csv":
			line, err = report.FormatRecordInCSV(c.parser.CsvTemplate, record)
		default:
			line, err = report.FormatRecordInCSV(c.parser.CsvTemplate, record)
		}

		if err != nil {
			return fmt.Errorf("failed to format the record: %v", err)
		}

		if _, err := file.Write(line); err != nil {
			return fmt.Errorf("failed to write to report: %v", err)
		}
	}

	return nil
}

func NewReportClient(opts ...report.Option) report.Client {
	options := report.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatal(err)
	}

	p := &report.Parser{}

	p.FilenameTemplate = p.ParseTemplate("filenameFormat", report.FilenameTemplate)
	p.CsvTemplate = p.ParseTemplate("csvFormat", report.CsvTemplate)

	c := &client{
		options: options,
		parser:  p,
	}

	return c
}
