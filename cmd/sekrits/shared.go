package main

import (
	"fmt"
	errors "github.com/hashicorp/errwrap"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
	"regexp"
	"strings"
)

func asFormat(as string, input string) string {
	if len(as) > 0 {
		return as
	}

	lookup := map[string]string{
		"yaml":  "yaml",
		"yml":   "yaml",
		"json":  "json",
		"json5": "json",
		"toml":  "toml",
	}

	var keys []string
	for ext := range lookup {
		keys = append(keys, ext)
	}

	pattern := fmt.Sprintf(".(%s)$", strings.Join(keys, "|"))

	re := regexp.MustCompile(pattern)
	fileExt := re.Find([]byte(strings.ToLower(input)))

	if len(fileExt) > 0 {
		return lookup[string(fileExt[1:])]
	}

	return "blob"
}

func processErrors(err error) error {
	ret := cli.MultiError{}
	errors.Walk(err, func(err error) {
		ret.Errors = append(ret.Errors, err)
	})
	return fmt.Errorf(ret.Error())
}

func getInputReader(input string) (io.Reader, error) {
	if input == "-" {
		return os.Stdin, nil
	}

	reader, err := os.Open(input)
	if err != nil {
		return reader, errors.Wrapf("Couldn't open file: {{err}}", err)
	}

	return reader, nil
}
