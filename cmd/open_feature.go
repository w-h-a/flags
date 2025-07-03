package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	openfeature "github.com/w-h-a/flags/internal/open_feature"
)

func OpenFeature(ctx *cli.Context) error {
	v, err := openfeature.Factory(
		ctx.String("baseUri"),
		ctx.Bool("insecure"),
		ctx.String("apiKey"),
		ctx.String("flag"),
		ctx.String("flagType"),
		"name",
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s evaluated as %+v\n", ctx.String("flag"), v)

	return nil
}
