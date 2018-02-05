package sekrits

import (
	"bytes"
	"encoding/base64"
	"fmt"
	errors "github.com/hashicorp/errwrap"
	"github.com/mikesimons/traverser"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

func (s *Sekrits) DecryptStructured(input io.Reader, format string) ([]byte, error) {
	codec, err := codecForFormat(format)
	if err != nil {
		return []byte(""), err
	}

	var inputData interface{}
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), fmt.Errorf("error reading %s input: %s", format, err.Error())
	}

	err = codec.Unmarshal(inputBytes, &inputData)
	if err != nil {
		return []byte(""), errors.Wrapf(fmt.Sprintf("error parsing %s", format), err)
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
					return traverser.Error(errors.Wrapf("Could not decrypt map value", err))
				}

				return traverser.Set(reflect.ValueOf(string(decrypted)))
			}
			return traverser.Noop()
		},
	}

	err = t.Traverse(reflect.ValueOf(inputData))

	if err != nil {
		return []byte(""), errors.Wrapf("error decrypting value", err)
	}

	decrypted, err := codec.Marshal(&inputData)
	if err != nil {
		return []byte(""), errors.Wrapf(fmt.Sprintf("error marshalling decrypted %s for output", format), err)
	}

	return decrypted, nil
}
