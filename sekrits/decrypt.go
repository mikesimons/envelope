package sekrits

import (
	"encoding/base64"
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

	decodedString, err := base64.StdEncoding.DecodeString(string(inputBytes))
	if err != nil {
		return []byte(""), err
	}

	decrypted, err := kr.Decrypt([]byte(decodedString))
	if err != nil {
		return []byte(""), err
	}

	return decrypted, nil
}
