package dynamodb

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/w-h-a/flags/internal/server/clients/writer"
	"github.com/w-h-a/flags/internal/server/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

const (
	table = "flags"
)

var AWSCFG aws.Config

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(config.Region()),
	)
	if err != nil {
		detail := "failed to register dynamodb writer with otel"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	AWSCFG = cfg
}

type client struct {
	options writer.Options
	conn    *dynamodb.Client
}

func (c *client) Write(ctx context.Context, key string, bs []byte) error {
	if _, err := c.conn.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item: map[string]types.AttributeValue{
				"Key":   &types.AttributeValueMemberS{Value: key},
				"Value": &types.AttributeValueMemberB{Value: bs},
			},
		},
	); err != nil {
		return err
	}

	return nil
}

func NewWriter(opts ...writer.Option) writer.Writer {
	options := writer.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		detail := "failed to validate dynamodb writer options"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	c := &client{
		options: options,
	}

	// just an http endpoint
	conn := dynamodb.NewFromConfig(
		AWSCFG,
		func(o *dynamodb.Options) {
			o.EndpointResolverV2 = &resolver{location: options.Location}
		},
	)

	c.conn = conn

	if _, err := c.conn.CreateTable(
		context.Background(),
		&dynamodb.CreateTableInput{
			TableName: aws.String(table),
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("Key"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Key"),
					KeyType:       types.KeyTypeHash,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(5),
			},
		},
	); err != nil && !strings.Contains(err.Error(), "ResourceInUseException") {
		detail := "failed to create table for dynamodb writer"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	return c
}
