package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/flags/cmd"
)

func main() {
	app := &cli.App{
		Name:  "flags",
		Usage: "flags client(s)/server",
		Commands: []*cli.Command{
			{
				Name: "server",
				Action: func(ctx *cli.Context) error {
					return cmd.Server(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("failed:", err)
	}
}
