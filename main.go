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
			{
				Name: "openfeature",
				Action: func(ctx *cli.Context) error {
					return cmd.OpenFeature(ctx)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "host",
						Usage:    "Provide hostname of server",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "port",
						Usage:    "Provide the port of the server",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "flag",
						Usage:    "Provide the flag key",
						Required: true,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("failed:", err)
	}
}
