package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mikesimons/envelope"
	"github.com/mikesimons/traverser"
	"gopkg.in/urfave/cli.v1"
)

const (
	errorInPlaceWithStdin = "You must provide a file to use --in-place"
	errorMissingKey       = "You must provide --key when setting values in json, yaml or toml"
)

func encryptBlob(file string, profile string, app *envelope.Envelope, c *cli.Context) error {
	secretReader, err := getInputReader(file)
	if err != nil {
		return err
	}

	if !c.Bool("no-trim") {
		secretReader = NewTrimReader(secretReader)
	}

	output, err := app.EncryptWithOpts(
		profile,
		secretReader,
		envelope.EncryptOpts{
			Encoder:    envelope.Base64Encoder,
			WithPrefix: c.Bool("with-prefix"),
		},
	)
	if err != nil {
		return err
	}

	return writeOutput(output, file, c.Bool("in-place"))
}

func encryptCommand() cli.Command {
	return cli.Command{
		Name:      "encrypt",
		Usage:     "Encrypt unencrypted data",
		ArgsUsage: "<file>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "profile",
				Value: "default",
				Usage: "Encryption profile name / id",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "Override type of input data (blob | yaml | json | toml)",
			},
			cli.BoolFlag{
				Name:  "no-trim",
				Usage: "Do not trim newlines from end of secret",
			},
			cli.BoolFlag{
				Name:  "in-place",
				Usage: "Write output back to input (requires file argument)",
			},
			cli.BoolFlag{
				Name:  "with-prefix",
				Usage: "Prefix encrypted string with envelope encryption marker (blob only)",
			},
			cli.StringFlag{
				Name:  "key",
				Usage: "Set encrypted value in this key (json, yaml & toml only)",
			},
		},
		Action: func(c *cli.Context) error {
			keyring := c.GlobalString("keyring")
			profile := c.String("profile")

			file := c.Args().Get(0)
			if file == "" {
				file = "-"
			}

			if file == "-" && c.Bool("in-place") {
				return processErrors(fmt.Errorf(errorInPlaceWithStdin))
			}

			app, err := envelope.WithYamlKeyring(keyring)
			if err != nil {
				return processErrors(err)
			}

			format := asFormat(c.String("format"), file)

			if format == "blob" {
				err := encryptBlob(file, profile, app, c)
				if err != nil {
					return processErrors(err)
				}
			} else {
				key := c.String("key")
				if key == "" {
					return processErrors(fmt.Errorf(errorMissingKey))
				}

				secretReader, err := getSecretReader(app, file, format, key, c.Bool("no-trim"))
				if err != nil {
					return processErrors(err)
				}

				fileReader, err := getInputReader(file)
				if err != nil {
					return processErrors(err)
				}

				output, err := app.InjectEncrypted(profile, fileReader, key, secretReader, format)
				if err != nil {
					return processErrors(err)
				}

				err = writeOutput(output, file, c.Bool("in-place"))
				if err != nil {
					return processErrors(err)
				}
			}

			return nil
		},
	}
}

func writeOutput(output []byte, file string, inPlace bool) error {
	var outputWriter io.WriteCloser
	var err error
	outputWriter = os.Stdout
	if inPlace {
		outputWriter, err = ioutil.TempFile(os.TempDir(), "envelope")
		if err != nil {
			return err
		}
		defer func() {
			os.Rename(outputWriter.(*os.File).Name(), file)
		}()
	}

	_, err = outputWriter.Write(output)
	outputWriter.Close()

	return err
}

func getSecretReader(app *envelope.Envelope, file string, format string, key string, noTrim bool) (io.Reader, error) {
	var secretReader io.Reader
	var err error

	stat, _ := os.Stdin.Stat()
	haveStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if haveStdin {
		secretReader, err = getInputReader("-")
		if err != nil {
			return nil, err
		}

		if !noTrim {
			secretReader = NewTrimReader(secretReader)
		}
	} else {
		fileReader, err := getInputReader(file)
		if err != nil {
			return nil, err
		}
		decrypted, err := app.DecryptStructuredAsMap(fileReader, format)
		tmp, err := traverser.GetKey(decrypted, strings.Split(key, "."))
		if err != nil {
			return nil, err
		}
		secretReader = strings.NewReader(fmt.Sprintf("%v", tmp))
	}

	return secretReader, nil
}
