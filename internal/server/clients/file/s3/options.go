package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/w-h-a/flags/internal/server/clients/file"
)

type awsConfigKey struct{}

func WithAWSConfig(cfg aws.Config) file.Option {
	return func(o *file.Options) {
		o.Context = context.WithValue(o.Context, awsConfigKey{}, cfg)
	}
}

func AWSConfig(ctx context.Context) (aws.Config, bool) {
	c, ok := ctx.Value(awsConfigKey{}).(aws.Config)
	return c, ok
}
