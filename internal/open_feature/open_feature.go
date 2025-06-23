package openfeature

import (
	"context"

	"github.com/open-feature/go-sdk-contrib/providers/ofrep"
	of "github.com/open-feature/go-sdk/openfeature"
)

func Factory(
	baseUri string,
	insecure bool,
	apiKey string,
	flagKey string,
	flagType string,
	name string,
) (any, error) {
	if insecure {
		baseUri = "http://" + baseUri
	} else {
		baseUri = "https://" + baseUri
	}

	provider := ofrep.NewProvider(
		baseUri,
		ofrep.WithBearerToken(apiKey),
	)

	if err := of.SetProviderAndWait(provider); err != nil {
		return nil, err
	}

	client := of.NewClient(name)

	evaluationCtx := of.NewEvaluationContext(
		"",
		map[string]any{},
	)

	var v any
	var err error

	switch flagType {
	case "int":
		v, err = client.IntValue(context.TODO(), flagKey, 0, evaluationCtx)
	case "float64":
		v, err = client.FloatValue(context.TODO(), flagKey, 0.0, evaluationCtx)
	case "string":
		v, err = client.StringValue(context.TODO(), flagKey, "", evaluationCtx)
	default:
		v, err = client.BooleanValue(context.TODO(), flagKey, false, evaluationCtx)
	}

	return v, err
}
