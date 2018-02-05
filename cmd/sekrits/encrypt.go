package main

import (
	"encoding/base64"
	"github.com/mikesimons/sekrits"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
)

func encryptCommand() cli.Command {
	return cli.Command{
		Name:  "encrypt",
		Usage: "Encrypt unencrypted data",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "with",
				Value: "default",
				Usage: "Encryption key alias / id",
			},
			cli.StringFlag{
				Name:  "for",
				Value: "yaml",
				Usage: "Type of data to encrypt for (blob | yaml | json | toml)",
			},
		},
		Action: func(c *cli.Context) error {
			keyring := c.GlobalString("keyring")
			alias := c.String("with")

			inputReader, err := getInputReader("-")
			if err != nil {
				return processErrors(err)
			}

			app, err := sekrits.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			var outputWriter io.WriteCloser
			var encrypted []byte

			outputWriter = base64.NewEncoder(base64.StdEncoding, os.Stdout)
			encrypted, err = app.Encrypt(alias, inputReader)
			if err != nil {
				return err
			}

			if c.String("for") != "blob" {
				os.Stdout.Write([]byte(app.Prefix))
			}

			outputWriter.Write(encrypted)
			outputWriter.Close()

			return nil
		},
	}
}
