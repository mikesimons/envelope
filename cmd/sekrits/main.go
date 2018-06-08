package main

import (
	"os"

	"github.com/ansel1/merry"
	"gopkg.in/urfave/cli.v1"
)

var COLLECT_DEBUG = true

func main() {
	app := cli.NewApp()
	app.Name = "sekrits"
	app.Usage = "Secrets encryption that doesn't suck"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "keyring",
			Value:  "keyring.yaml",
			Usage:  "Keyring file / url",
			EnvVar: "SEKRITS_KEYRING",
		},
	}

	app.Commands = []cli.Command{
		addKeyCommand(),
		encryptCommand(),
		decryptCommand(),
	}

	merry.SetStackCaptureEnabled(COLLECT_DEBUG)

	app.Run(os.Args)
}
