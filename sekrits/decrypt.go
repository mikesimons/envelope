package sekrits

import (
	"fmt"
	errors "github.com/hashicorp/errwrap"
	"io"
	"io/ioutil"
)

func (s *Sekrits) Decrypt(input io.Reader) ([]byte, error) {
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), fmt.Errorf("error reading value: %s", err.Error())
	}

	decrypted, err := s.Keyring.Decrypt(inputBytes)
	if err != nil {
		return []byte(""), errors.Wrapf("error decrypting input", err)
	}

	return decrypted, nil
}
