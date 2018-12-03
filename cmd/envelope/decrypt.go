package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"

	"log"

	"github.com/mikesimons/envelope"
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
				Usage: "Type of input data (blob | yaml | json | toml)",
			},
			cli.StringFlag{
				Name:  "on-error",
				Usage: "Method for handling decryption errors (unset | replace:<value> | ignore | exit)",
				Value: "exit",
			},
		},
		Action: func(c *cli.Context) error {
			keyring := c.GlobalString("keyring")

			input := c.Args().Get(0)
			if input == "" {
				input = "-"
			}

			var err error
			outputWriter := os.Stdout

			inputReader, err := getInputReader(input)
			if err != nil {
				return processErrors(err)
			}

			app, err := envelope.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			var decrypted []byte
			as := asFormat(c.String("format"), input)

			switch as {
			case "blob":
				decrypted, err = app.Decrypt(base64.NewDecoder(base64.StdEncoding, inputReader))
				if err != nil {
					return processErrors(err)
				}
			default:
				handler, err := structuredErrorHandler(c.String("on-error"))
				if err != nil {
					return processErrors(err)
				}
				app.StructuredErrorBehaviour = handler
				decrypted, err = app.DecryptStructured(inputReader, as)
				if err != nil {
					return processErrors(err)
				}
			}

			outputWriter.Write(decrypted)
			outputWriter.Close()

			return nil
		},
	}
}

func structuredErrorHandler(strategy string) (func(error) (traverser.Op, error), error) {
	values := strings.SplitN(strategy, ":", 2)

	switch values[0] {
	case "unset":
		return func(e error) (traverser.Op, error) {
			// todo debug log e
			return traverser.Unset()
		}, nil

	case "replace":
		replace := values[0]
		if replace == "" {
			replace = "ERROR"
		}

		return func(e error) (traverser.Op, error) {
			// todo debug log e
			return traverser.Set(reflect.ValueOf(replace))
		}, nil

	case "ignore":
		return func(e error) (traverser.Op, error) {
			// todo debug log e
			return traverser.Noop()
		}, nil

	case "exit":
		return func(e error) (traverser.Op, error) {
			// todo debug log e
			log.Print(e)
			return traverser.ErrorNoop(e)
		}, nil
	}

	return nil, fmt.Errorf("Invalid error handling stategy '%s'. Valid values: unset, replace:<value>, ignore, exit", strategy)
}
