package envelope

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ansel1/merry"
	"github.com/mikesimons/traverser"
)

// InjectEncrypted inject the secret encrypted.
func (s *Envelope) InjectEncrypted(alias string, input io.Reader, key string, value io.Reader, format string) ([]byte, error) {
	inputData, codec, err := injectHelper(input, format)
	if err != nil {
		return nil, err
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

	err = setKey(inputData, encrypted, key)
	if err != nil {
		return nil, err
	}

	ret, err := codec.Marshal(&inputData)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("error marshalling output")
	}

	return ret, nil
}

// InjectNotEncrypted inject the secret not encrypted. This can be used together with Encrypt()
func (s *Envelope) InjectNotEncrypted(alias string, input io.Reader, key string, value []byte, format string) ([]byte, error) {
	inputData, codec, err := injectHelper(input, format)
	if err != nil {
		return nil, err
	}

	encrypted := value

	err = setKey(inputData, encrypted, key)
	if err != nil {
		return nil, err
	}

	ret, err := codec.Marshal(&inputData)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("error marshalling output")
	}

	return ret, nil
}

// injectHelper is a helper function for InjectEncrypted and InjectNotEncrypted
func injectHelper(input io.Reader, format string) (interface{}, structuredCodec, error) {
	codec, err := codecForFormat(format)
	if err != nil {
		return nil, structuredCodec{}, merry.Wrap(err).WithUserMessage("unrecognized format").WithValue("format", format)
	}

	var inputData interface{}
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, structuredCodec{}, err
	}

	err = codec.Unmarshal(inputBytes, &inputData)
	if err != nil {
		return nil, structuredCodec{}, merry.Wrap(err).WithUserMessage("could not decode input").WithValue("format", format)
	}

	return inputData, codec, nil
}

func setKey(inputData interface{}, encrypted interface{}, key string) error {
	splitKey := strings.Split(key, ".")

	err := traverser.SetKey(inputData, splitKey, fmt.Sprintf("%s", encrypted))
	if err != nil {
		return merry.Wrap(err).WithValue("key", key)
	}

	return nil
}
