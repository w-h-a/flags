package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

func Flags(ctx *cli.Context) error {
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix("flags"),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

	filePath := ctx.String("filePath")

	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	bs, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fs, err := flags.Factory(
		bs,
		ctx.String("format"),
	)
	if err != nil {
		return err
	}

	for k := range fs {
		log.Infof("flag %s is valid", k)
	}

	return nil
}
