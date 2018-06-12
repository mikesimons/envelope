package sekrits

import (
	"io"
	"io/ioutil"
)

func (s *Sekrits) Decrypt(input io.Reader) ([]byte, error) {
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), err
	}

	decrypted, err := s.Keyring.Decrypt(inputBytes)
	if err != nil {
		return []byte(""), err
	}

	return decrypted, nil
}
