package s3

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/w-h-a/flags/internal/server/clients/report"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options report.Options
	awsCfg  aws.Config
	parser  *report.Parser
}

func (c *client) Create(ctx context.Context, records []report.Record) error {
	dir, err := os.MkdirTemp("", "flags_s3_report")
	if err != nil {
		return err
	}

	defer os.Remove(dir)

	sess, err := session.NewSession(&c.awsCfg)
	if err != nil {
		return err
	}

	s3Uploader := s3manager.NewUploader(sess)

	filename, err := c.parser.ParseFilename(c.options.Format)
	if err != nil {
		return err
	}

	filePath := dir + "/" + filename

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

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

	file.Close()

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		opened, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			return err
		}

		if _, err := s3Uploader.UploadWithContext(
			ctx,
			&s3manager.UploadInput{
				Bucket: aws.String(c.options.Dir),
				Key:    aws.String(file.Name()),
				Body:   opened,
			},
		); err != nil {
			return err
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

	if cfg, ok := AWSConfig(options.Context); !ok {
		log.Fatalf("aws config is required for s3 retriever")
	} else {
		c.awsCfg = cfg
	}

	return c
}
