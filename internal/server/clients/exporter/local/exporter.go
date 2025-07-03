package local

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/w-h-a/flags/internal/server/clients/exporter"
)

type client struct {
	options exporter.Options
	parser  *exporter.Parser
}

func (c *client) Export(ctx context.Context, records []exporter.Record) error {
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
			line, err = exporter.FormatRecordInCSV(c.parser.CsvTemplate, record)
		default:
			line, err = exporter.FormatRecordInCSV(c.parser.CsvTemplate, record)
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

func NewExporter(opts ...exporter.Option) exporter.Exporter {
	options := exporter.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		slog.ErrorContext(context.Background(), "failed to validate local exporter options", "error", err)
		os.Exit(1)
	}

	p := &exporter.Parser{}

	p.FilenameTemplate = p.ParseTemplate("filenameFormat", exporter.FilenameTemplate)
	p.CsvTemplate = p.ParseTemplate("csvFormat", exporter.CsvTemplate)

	c := &client{
		options: options,
		parser:  p,
	}

	return c
}
