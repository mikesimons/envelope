package main

import (
	"fmt"
	"net/url"

	"github.com/ansel1/merry"
	"github.com/mikesimons/envelope"
	"gopkg.in/urfave/cli.v1"
)

// profileAddCommand implements the `profile add` command
// e.g $ envelope profile add test kms://arn
func profileAddCommand() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "Add a profile to the keyring.",
		Description: `Master key DSN is of the form <provider>://<dsn>.
   The only supported provider is currently awskms.

   Example:

     envelope profile add test \
       awskms://arn:aws:kms:us-east-1:111111111111:key/abcdef01-2345-6789-abcd-ef0123456789 \
       --context='role=somerole&test=test'`,
		ArgsUsage: "<alias> <master key dsn>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "context",
				Usage: "Encryption context in URL param format (e.g. key1=value1&key2=value2)",
			},
		},
		Action: func(c *cli.Context) error {
			keyring := c.GlobalString("keyring")
			alias := c.Args().Get(0)
			if alias == "" {
				cli.ShowCommandHelp(c, "add")
				return cli.NewExitError("Missing alias argument", 1)
			}

			providerDsn := c.Args().Get(1)
			if providerDsn == "" {
				cli.ShowCommandHelp(c, "add")
				return cli.NewExitError("Missing master key dsn argument", 1)
			}

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
