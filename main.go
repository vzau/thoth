package main

import (
	"fmt"
	"os"

	"github.com/dhawton/log4g"
	"github.com/urfave/cli/v2"
	"github.com/vzau/thoth/internal/server"
	"github.com/vzau/thoth/pkg/version"
)

var log = log4g.Category("main")

func main() {
	app := &cli.App{
		Name:                 "thoth",
		Usage:                "ZAU API",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Set log level to debug",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Start API Server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Set port to listen on",
						Value:   3000,
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("debug") {
						log4g.SetLogLevel(log4g.DEBUG)
					}
					log.Info("Starting server")
					server.Run(c.Int("port"))
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version",
				Action: func(c *cli.Context) error {
					fmt.Printf("Version %s\n", version.FriendlyVersion())
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
