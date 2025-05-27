package dynamodb

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	transport "github.com/aws/smithy-go/endpoints"
)

type resolver struct {
	location string
}

func (r *resolver) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (transport.Endpoint, error) {
	u, err := url.Parse(r.location)
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}
