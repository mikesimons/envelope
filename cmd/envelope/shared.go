package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ansel1/merry"
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
	var lines []string
	e := merry.Wrap(err)

	if msg := merry.UserMessage(e); msg != "" {
		lines = append(lines, fmt.Sprintf("%s\n", msg))
	}

	if msg := merry.Message(e); msg != "" {
		lines = append(lines, fmt.Sprintf("%s\n", msg))
	}

	if COLLECT_DEBUG {
		lines = append(lines, merry.Stacktrace(e))
	}

	var extra []string
	for k, v := range merry.Values(e) {
		k = fmt.Sprintf("%v", k)
		if k == "stack" || k == "message" || k == "user message" {
			continue
		}
		extra = append(extra, fmt.Sprintf("%v: %v", k, v))
	}
	lines = append(lines, strings.Join(extra, ", "))

	return fmt.Errorf(strings.Join(lines, "\n"))
}

func getInputReader(input string) (io.Reader, error) {
	if input == "-" {
		return os.Stdin, nil
	}

	reader, err := os.Open(input)
	if err != nil {
		return reader, merry.Wrap(err).WithValue("input", input)
	}

	return reader, nil
}

type TrimReader struct {
	r io.Reader
}

func NewTrimReader(r io.Reader) io.Reader {
	return &TrimReader{r: r}
}

func (t *TrimReader) Read(p []byte) (n int, err error) {
	count, err := t.r.Read(p)
	if count < len(p) && err == nil {
		pos := count
		for pos >= 0 {
			eval := p[:pos][len(p[:pos])-1]
			if eval == byte(0) || eval == byte('\n') || eval == byte('\r') {
				pos = pos - 1
				continue
			} else {
				break
			}
		}

		for i := pos; i < count; i++ {
			p[i] = byte(0)
		}

		return pos, err
	}
	return count, err
}
