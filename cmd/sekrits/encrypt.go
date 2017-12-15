package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mikesimons/sekrits/sekrits"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
)

func encryptCommand() cli.Command {
	return cli.Command{
		Name:  "encrypt",
		Usage: "Encrypt unencrypted data",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowCommandHelp(c, "encrypt")
				fmt.Println("")
				return cli.NewExitError("Error: Not enough arguments", 1)
			}

			keyring := c.GlobalString("keyring")
			alias := c.Args().Get(0)
			input := c.Args().Get(1)

			var inputReader io.Reader
			outputWriter := base64.NewEncoder(base64.StdEncoding, os.Stdout)

			if input == "-" {
				inputReader = os.Stdin
			} else {
				inputReader, _ = os.Open(input)
			}

			encrypted, err := sekrits.Encrypt(keyring, alias, inputReader)
			if err != nil {
				return err
			}

			outputWriter.Write(encrypted)

			return nil
		},
	}
}
