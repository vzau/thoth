/*
ZAU Thoth API
Copyright (C) 2021 Daniel A. Hawton (daniel@hawton.org)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
