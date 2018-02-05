package sekrits

import (
	"fmt"
	"io"
	"io/ioutil"
)

func (s *Sekrits) Encrypt(alias string, input io.Reader) ([]byte, error) {
	key, ok := s.Keyring.GetKey(alias)
	if !ok {
		return []byte(""), fmt.Errorf("Couldn't find key with alias or id '%s'", alias)
	}

	data, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), err
	}

	encrypted, err := key.EncryptBytes(data)
	if err != nil {
		return []byte(""), err
	}

	return encrypted, nil
}
