package envelope

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/ansel1/merry"
	"github.com/mikesimons/traverser"
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
	err = traverser.SetKey(inputData, splitKey, string(encrypted))
	if err != nil {
		return []byte(""), merry.Wrap(err).WithValue("key", key)
	}

	ret, err := codec.Marshal(&inputData)
	if err != nil {
		return []byte(""), merry.Wrap(err).WithUserMessage("error marshalling output")
	}

	return ret, nil
}
