package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"

	"log"

	"github.com/mikesimons/sekrits"
	"github.com/mikesimons/traverser"
	"gopkg.in/urfave/cli.v1"
)

func decryptCommand() cli.Command {
	return cli.Command{
		Name:  "decrypt",
		Usage: "Decrypt encrypted data",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "format",
				Usage: "Type of data to encrypt (blob | yaml | json | toml)",
			},
			cli.BoolFlag{
				Name:  "unset-errors",
				Usage: "Unset keys if they can't be decrypted (only applies to structured decryption)",
			},
			cli.StringFlag{
				Name:  "default-error-value",
				Usage: "Set keys to this value if they can't be decrypted (only applies to structured decryption)",
				Value: "ERROR",
			},
			cli.BoolFlag{
				Name:  "ignore-errors",
				Usage: "Ignore decryption errors and leave the encrypted string in place (only applies to structured decryption)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowCommandHelp(c, "decrypt")
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
			as := asFormat(c.String("format"), input)

			switch as {
			case "blob":
				decrypted, err = app.Decrypt(base64.NewDecoder(base64.StdEncoding, inputReader))
			default:
				app.StructuredErrorBehaviour = structuredErrorHandler(c)
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

func structuredErrorHandler(c *cli.Context) func(error) (traverser.Op, error) {
	if c.Bool("ignore-errors") {
		return func(e error) (traverser.Op, error) {
			log.Print(e)
			return traverser.Noop()
		}
	}

	if c.Bool("unset-errors") {
		return func(e error) (traverser.Op, error) {
			log.Print(e)
			return traverser.Unset()
		}
	}

	return func(e error) (traverser.Op, error) {
		log.Print(e)
		return traverser.Set(reflect.ValueOf(c.String("default-error-value")))
	}
}
