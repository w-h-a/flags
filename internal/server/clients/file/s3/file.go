package s3

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options file.Options
	cfg     aws.Config
	parser  *file.Parser
}

func (c *client) Read(ctx context.Context) (map[string]*file.Flag, error) {
	f, err := os.CreateTemp("", "feature_flag")
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&c.cfg)
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)

	file := ""

	// TODO: generalize
	if len(c.options.Files) > 0 {
		file = c.options.Files[0]
	}

	req := &s3.GetObjectInput{
		Bucket: aws.String(c.options.Dir),
		Key:    aws.String(file),
	}

	if _, err := downloader.DownloadWithContext(ctx, f, req); err != nil {
		return nil, fmt.Errorf("failed to download file from s3 %q: %v", file, err)
	}

	bs, err := os.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	return c.parser.ParseFlags(bs)
}

func NewFileClient(opts ...file.Option) file.Client {
	options := file.NewOptions(opts...)

	c := &client{
		options: options,
		parser:  &file.Parser{},
	}

	if cfg, ok := AWSConfig(options.Context); !ok {
		log.Fatalf("aws config is required for s3 retriever")
		return nil
	} else {
		c.cfg = cfg
	}

	return c
}
