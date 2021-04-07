package main

import (
	"os"

	"github.com/jalavosus/rsync-daemon/internal"
	"github.com/jalavosus/rsync-daemon/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "rsync-daemon",
		Authors: []*cli.Author{
			{
				Name:  "jalavosus",
				Email: "alavosus.james@gmail.com",
			},
		},
		Version: internal.AppVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage: "Path to a config yaml to override default settings. \n" +
					"If a directory is passed, rsync-daemon assumes that a file with the " +
					"default config filename (" + config.DefaultConfigFilename + ") exists in that directory. \n" +
					"If a full path to a file is passed, that specific file will be used. \n" +
					"If a full file path is passed, the config file can be in either JSON or YAML format, and must " +
					"include the extension .json or .yaml",
				Value:    config.DefaultConfigPath,
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "ignore-validation",
				Usage:    "If true, all config validation errors will be ignored. Not recommended.",
				Required: false,
				Value:    false,
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
