package main

import (
	"fmt"
	"net/url"

	"github.com/ansel1/merry"
	"github.com/mikesimons/envelope"
	"gopkg.in/urfave/cli.v1"
)

// addKeyCommand implements the add-key command
// e.g $ envelope add-key test kms://arn
func addKeyCommand() cli.Command {
	return cli.Command{
		Name:  "addkey",
		Usage: "Add a key to the keyring.",
		Description: `Master key DSN is of the form <provider>://<dsn>. The only supported provider is currently awskms.
   Example: envelope add-key mytestkey awskms://arn:aws:kms:us-east-1:111111111111:key/abcdef01-2345-6789-abcd-ef0123456789`,
		ArgsUsage: "<alias> <master key dsn>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "context",
				Usage: "Encryption context in URL param format (e.g. key1=value1&key2=value2)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowCommandHelp(c, "addkey")
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

			app, err := envelope.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			keyID, err := app.AddKey(alias, providerDsn, context)
			if err != nil {
				return processErrors(err)
			}

			fmt.Printf("Added key %s (%s)", keyID, alias)

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
		return nil, merry.Wrap(err).WithMessage("Unable to parse encryption context")
	}

	output := make(map[string]string)
	for k, vs := range values {
		output[k] = vs[0]
	}

	return output, nil
}
