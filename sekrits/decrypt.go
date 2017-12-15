package sekrits

import (
	"github.com/mikesimons/sekrits/keyring"
	"io"
	"io/ioutil"
)

func Decrypt(keyringPath string, input io.Reader) ([]byte, error) {
	kr, err := keyring.Load(keyringPath)
	if err != nil {
		return []byte(""), err
	}

	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), err
	}

	decrypted, err := kr.Decrypt(inputBytes)
	if err != nil {
		return []byte(""), err
	}

	return decrypted, nil
}
