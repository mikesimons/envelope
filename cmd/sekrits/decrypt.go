package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mikesimons/sekrits/sekrits"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func decryptCommand() cli.Command {
	return cli.Command{
		Name:  "decrypt",
		Usage: "Decrypt encrypted data",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "as",
				Usage: "Type of data to encrypt (blob | yaml | json | toml)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowCommandHelp(c, "encrypt")
				fmt.Println("")
				return cli.NewExitError("Error: Not enough arguments", 1)
			}

			keyring := c.GlobalString("keyring")
			input := c.Args().Get(0)

			var err error
			outputWriter := os.Stdout

			inputReader, err := getInputReader(input)
			if err != nil {
				return processErrors(err)
			}

			app, err := sekrits.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			var decrypted []byte
			as := asFormat(c.String("as"), input)

			switch as {
			case "blob":
				// TODO optionally read app.Prefix
				decrypted, err = app.Decrypt(base64.NewDecoder(base64.StdEncoding, inputReader))
			default:
				decrypted, err = app.DecryptStructured(inputReader, as)
			}

			if err != nil {
				return processErrors(err)
			}

			outputWriter.Write(decrypted)
			outputWriter.Close()

			return nil
		},
	}
}
