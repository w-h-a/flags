package openfeature

import (
	"context"
	"fmt"

	"github.com/open-feature/go-sdk-contrib/providers/ofrep"
	of "github.com/open-feature/go-sdk/openfeature"
)

func Factory(host string, port int, flagKey, name string) (any, error) {
	provider := ofrep.NewProvider(
		fmt.Sprintf("http://%s:%d", host, port),
	)

	if err := of.SetProviderAndWait(provider); err != nil {
		return nil, err
	}

	client := of.NewClient(name)

	evaluationCtx := of.NewEvaluationContext(
		"",
		map[string]any{},
	)

	v, err := client.BooleanValue(context.TODO(), flagKey, false, evaluationCtx)

	return v, err
}
