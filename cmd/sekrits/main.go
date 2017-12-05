package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "sekrits"
	app.Usage = "Secrets management that doesn't suck"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "keyring",
			Value: "keyring.yaml",
			Usage: "Keyring file / url",
			EnvVar: "SEKRITS_KEYRING",
		},
	}

	app.Commands = []cli.Command{
		addKeyCommand(),
	}

	app.Run(os.Args)
}
