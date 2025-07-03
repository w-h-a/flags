package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/flags/internal/flags"
)

func Flags(ctx *cli.Context) error {
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
		fmt.Printf("flag %s is valid\n", k)
	}

	return nil
}
