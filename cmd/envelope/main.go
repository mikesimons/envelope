package main

import (
	"os"

	"github.com/mikesimons/envelope/keysvc"

	"github.com/ansel1/merry"
	"github.com/mikesimons/envelope/keysvc/awskms"
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
		cli.StringFlag{
			Name:  "aws-region",
			Usage: "AWS region",
		},
		cli.StringFlag{
			Name:  "aws-role",
			Usage: "AWS role ARN",
		},
		cli.StringFlag{
			Name:  "aws-profile",
			Usage: "AWS profile name",
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

		sess := awsSession(c.String("aws-profile"), c.String("aws-region"), c.String("aws-role"))
		keysvc.AddKeyServiceFn("awskms", func() (keysvc.KeyServiceProvider, error) {
			return awskms.New(sess)
		})

		return nil
	}

	app.Run(os.Args)
}
