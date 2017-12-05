package main

import (
	"fmt"
	"github.com/mikesimons/sekrits/sekrits"
	"gopkg.in/urfave/cli.v1"
)

// addKeyCommand implements the add-key command
// e.g $ sekrits add-key test kms://arn
func addKeyCommand() cli.Command {
	return cli.Command{
		Name:      "add-key",
		Usage:     "Add a key to the keyring",
		ArgsUsage: "<alias> <master key dsn>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowCommandHelp(c, "add-key")
				fmt.Println("")
				return cli.NewExitError("Error: Not enough arguments", 1)
			}

			keyring := c.GlobalString("keyring")
			alias := c.Args().Get(0)
			providerDsn := c.Args().Get(1)

			keyId, err := sekrits.AddKey(keyring, alias, providerDsn)
			if err != nil {
				return err
			}

			fmt.Printf("Added key %s (%s)", keyId, alias)

			return nil
		},
	}
}
