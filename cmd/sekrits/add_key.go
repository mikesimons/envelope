package main

import (
	"fmt"
	"net/url"

	errors "github.com/hashicorp/errwrap"
	"github.com/mikesimons/sekrits"
	"gopkg.in/urfave/cli.v1"
)

// addKeyCommand implements the add-key command
// e.g $ sekrits add-key test kms://arn
func addKeyCommand() cli.Command {
	return cli.Command{
		Name:      "add-key",
		Usage:     "Add a key to the keyring",
		ArgsUsage: "<alias> <master key dsn>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "context",
				Usage: "Encryption context in URL param format (e.g. key1=value1&key2=value2)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowCommandHelp(c, "add-key")
				fmt.Println("")
				return cli.NewExitError("Error: Not enough arguments", 1)
			}

			keyring := c.GlobalString("keyring")
			alias := c.Args().Get(0)
			providerDsn := c.Args().Get(1)
			context, err := parseEncryptionContext(c.String("context"))
			if err != nil {
				return processErrors(err)
			}

			app, err := sekrits.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			keyId, err := app.AddKey(alias, providerDsn, context)
			if err != nil {
				return processErrors(err)
			}

			fmt.Printf("Added key %s (%s)", keyId, alias)

			return nil
		},
	}
}

func parseEncryptionContext(input string) (map[string]string, error) {
	if input == "" {
		return nil, nil
	}

	values, err := url.ParseQuery(input)
	if err != nil {
		return nil, errors.Wrapf("Unable to parse encryption context", err)
	}

	output := make(map[string]string)
	for k, vs := range values {
		output[k] = vs[0]
	}

	return output, nil
}
