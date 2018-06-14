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
				Name:  "key",
				Value: "default",
				Usage: "Encryption key alias / id",
			},
			cli.BoolFlag{
				Name:  "blob",
				Usage: "Encrypt as blob",
			},
			cli.BoolFlag{
				Name:  "notrim",
				Usage: "Do not trim newlines from end of input",
			},
		},
		Action: func(c *cli.Context) error {
			keyring := c.GlobalString("keyring")
			alias := c.String("key")

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
			var encrypted []byte

			outputWriter = base64.NewEncoder(base64.StdEncoding, os.Stdout)
			encrypted, err = app.Encrypt(alias, inputReader)
			if err != nil {
				return processErrors(err)
			}

			if !c.Bool("blob") {
				os.Stdout.Write([]byte(app.Prefix))
			}

			outputWriter.Write(encrypted)
			outputWriter.Close()

			return nil
		},
	}
}
