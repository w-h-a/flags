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
					&cli.BoolFlag{
						Name:  "insecure",
						Usage: "Provide if using http",
					},
					&cli.StringFlag{
						Name:     "apiKey",
						Usage:    "Provide the api key for the server",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "flag",
						Usage:    "Provide the flag key",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "flagType",
						Usage:    "Provide the flag type (i.e., bool, int, float64, string)",
						Required: true,
					},
				},
			},
			{
				Name: "flags",
				Action: func(ctx *cli.Context) error {
					return cmd.Flags(ctx)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "filePath",
						Usage:    "Provide the file path to the flags",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "format",
						Usage:    "Provide the format of the flags (yaml or json)",
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
