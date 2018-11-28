package main

import (
	"encoding/base64"
	"io"
	"os"

	"github.com/mikesimons/envelope"
	"gopkg.in/urfave/cli.v1"
)

func encryptCommand() cli.Command {
	return cli.Command{
		Name:      "encrypt",
		Usage:     "Encrypt unencrypted data",
		ArgsUsage: "<file>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "profile",
				Value: "default",
				Usage: "Encryption profile name / id",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "Type of input data (blob | yaml | json | toml)",
			},
			cli.BoolFlag{
				Name:  "notrim",
				Usage: "Do not trim newlines from end of input",
			},
			cli.StringFlag{
				Name:  "key",
				Value: "Key to encrypt (structured formats only)",
			},
			cli.StringFlag{
				Name:  "secret",
				Value: "Secret",
			},
		},
		Action: func(c *cli.Context) error {
			keyring := c.GlobalString("keyring")
			alias := c.String("profile")

			file := c.Args().Get(0)
			if file == "" {
				file = "-"
			}

			inputReader, err := getInputReader(file)
			if err != nil {
				return processErrors(err)
			}

			if !c.Bool("notrim") {
				inputReader = NewTrimReader(inputReader)
			}

			app, err := envelope.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			var outputWriter io.WriteCloser
			var output []byte

			if c.String("key") != "" {
				as := asFormat(c.String("format"), file)
				outputWriter = os.Stdout
				output, err = app.EncryptInPlace(inputReader, c.String("key"), c.String("secret"), as)
				if err != nil {
					return processErrors(err)
				}

			} else {
				outputWriter = base64.NewEncoder(base64.StdEncoding, os.Stdout)
				output, err = app.Encrypt(alias, inputReader)
				if err != nil {
					return processErrors(err)
				}

				if !c.Bool("blob") {
					os.Stdout.Write([]byte(app.Prefix))
				}

			}

			outputWriter.Write(output)
			outputWriter.Close()

			return nil
		},
	}
}
