package envelope

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/ansel1/merry"
	"github.com/mikesimons/traverser"
)

// DecryptStructured will parse the input as format and use the encryption prefix to automatically identify and decrypt encrypted values
func (s *Envelope) DecryptStructured(input io.Reader, format string) ([]byte, error) {
	codec, err := codecForFormat(format)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("unrecognized format").WithValue("format", format)
	}

	var inputData interface{}
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), err
	}

	err = codec.Unmarshal(inputBytes, &inputData)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("could not decode input").WithValue("format", format)
	}

	t := &traverser.Traverser{
		Node: func(keys []string, val reflect.Value) (traverser.Op, error) {
			data := val.Interface()
			str, ok := data.(string)
			if ok && strings.HasPrefix(str, s.Prefix) {
				v := str[len(s.Prefix):]
				bytesReader := bytes.NewReader([]byte(v))
				inputReader := base64.NewDecoder(base64.StdEncoding, bytesReader)
				decrypted, err := s.Decrypt(inputReader)

				if err != nil {
					return s.StructuredErrorBehaviour(merry.Wrap(err).WithValue("key", strings.Join(keys, ".")))
				}

				return traverser.Set(reflect.ValueOf(string(decrypted)))
			}
			return traverser.Noop()
		},
	}

	err = t.Traverse(reflect.ValueOf(inputData))

	if err != nil {
		return []byte(""), err
	}

	decrypted, err := codec.Marshal(&inputData)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("error marshalling output")
	}

	return decrypted, nil
}
