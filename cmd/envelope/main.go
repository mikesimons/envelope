package main

import (
	"os"

	"github.com/ansel1/merry"
	"gopkg.in/urfave/cli.v1"
)

var COLLECT_DEBUG = true
var envelope_version_string = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "envelope"
	app.Usage = "Envelope secrets encryption"
	app.Version = envelope_version_string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "keyring",
			Value:  "keyring.yaml",
			Usage:  "Keyring file / url",
			EnvVar: "ENVELOPE_KEYRING",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "profile",
			Usage: "Profile related commands",
			Subcommands: []cli.Command{
				profileAddCommand(),
			},
		},
		encryptCommand(),
		decryptCommand(),
	}

	app.Before = func(c *cli.Context) error {
		COLLECT_DEBUG = c.Bool("debug")
		merry.SetStackCaptureEnabled(c.Bool("debug"))
		return nil
	}

	app.Run(os.Args)
}
