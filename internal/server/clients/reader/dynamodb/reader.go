package dynamodb

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/pkg/telemetry/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

const (
	table = "flags"
)

var AWSCFG aws.Config

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		// TODO: grab from config
		awsconfig.WithRegion("local"),
	)
	if err != nil {
		log.Fatal(err)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	AWSCFG = cfg
}

type client struct {
	options reader.Options
	conn    *dynamodb.Client
}

func (c *client) ReadByKey(ctx context.Context, key string) ([]byte, error) {
	rsp, err := c.conn.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(table),
			Key: map[string]types.AttributeValue{
				"Key": &types.AttributeValueMemberS{Value: key},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if rsp.Item == nil {
		return []byte{}, reader.ErrRecordNotFound
	}

	data, ok := rsp.Item["Value"].(*types.AttributeValueMemberB)
	if !ok {
		return nil, reader.ErrRecordNotFound
	}

	return data.Value, nil
}

func (c *client) Read(ctx context.Context) ([]byte, error) {
	paginator := dynamodb.NewScanPaginator(
		c.conn,
		&dynamodb.ScanInput{
			TableName: aws.String(table),
		},
	)

	records := []*reader.Record{}

	for paginator.HasMorePages() {
		rs := []*reader.Record{}

		rsp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if err := attributevalue.UnmarshalListOfMaps(rsp.Items, &rs); err != nil {
			return nil, err
		}

		records = append(records, rs...)
	}

	result := []byte{}

	for _, record := range records {
		result = append(result, []byte("\n")...)
		result = append(result, record.Value...)
	}

	return result, nil
}

func NewReader(opts ...reader.Option) reader.Reader {
	options := reader.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return c
}
