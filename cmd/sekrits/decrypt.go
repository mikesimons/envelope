package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mikesimons/sekrits/sekrits"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
)

func decryptCommand() cli.Command {
	return cli.Command{
		Name:  "decrypt",
		Usage: "Decrypt encrypted data",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowCommandHelp(c, "encrypt")
				fmt.Println("")
				return cli.NewExitError("Error: Not enough arguments", 1)
			}

			keyring := c.GlobalString("keyring")
			input := c.Args().Get(0)

			var inputReader io.Reader
			outputWriter := os.Stdout

			if input == "-" {
				inputReader = os.Stdin
			} else {
				inputReader, _ = os.Open(input)
			}

			inputReader = base64.NewDecoder(base64.StdEncoding, inputReader)

			decrypted, err := sekrits.Decrypt(keyring, inputReader)
			if err != nil {
				return err
			}

			outputWriter.Write(decrypted)

			return nil
		},
	}
}
