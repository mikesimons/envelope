package envelope

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/ansel1/merry"
)

func (s *Envelope) InjectEncrypted(alias string, input io.Reader, key string, value io.Reader, format string) ([]byte, error) {
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

	encrypted, err := s.EncryptWithOpts(
		alias,
		value,
		EncryptOpts{
			Encoder:    Base64Encoder,
			WithPrefix: true,
		},
	)

	if err != nil {
		return []byte(""), err
	}

	splitKey := strings.Split(key, ".")
	err = set(reflect.ValueOf(inputData), splitKey, []string{}, reflect.ValueOf(string(encrypted)))
	if err != nil {
		return []byte(""), merry.Wrap(err).WithValue("key", key)
	}

	ret, err := codec.Marshal(&inputData)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("error marshalling output")
	}

	return ret, nil
}

func set(data reflect.Value, path []string, traversed []string, value reflect.Value) error {
	nextKey := path[0]
	nextPath := path[1:]
	var zeroVal reflect.Value

	switch data.Kind() {
	case reflect.Interface:
		return set(data.Elem(), path, traversed, value)
	case reflect.Map:
		nextKeyVal := reflect.ValueOf(nextKey)
		nextVal := data.MapIndex(nextKeyVal)

		if len(nextPath) == 0 {
			data.SetMapIndex(nextKeyVal, value)
			return nil
		}

		if nextVal == zeroVal {
			nextVal = reflect.ValueOf(make(map[string]interface{}))
			data.SetMapIndex(nextKeyVal, nextVal)
		}

		traversed = append(traversed, nextKey)
		return set(nextVal, nextPath, traversed, value)
	default:
		return fmt.Errorf("Can't set key because %s is a %s", strings.Join(traversed, "."), data.Kind().String())
	}
}
