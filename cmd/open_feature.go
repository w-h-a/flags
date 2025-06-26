package cmd

import (
	"github.com/urfave/cli/v2"
	openfeature "github.com/w-h-a/flags/internal/open_feature"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

func OpenFeature(ctx *cli.Context) error {
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix("openfeature"),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

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

	log.Infof("%s evaluated as %+v", ctx.String("flag"), v)

	return nil
}
