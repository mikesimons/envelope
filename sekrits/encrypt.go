package sekrits

import (
	"encoding/base64"
	"fmt"
	"github.com/mikesimons/sekrits/keyring"
	"io"
	"io/ioutil"
)

func Encrypt(keyringPath string, alias string, input io.Reader) ([]byte, error) {
	kr, err := keyring.Load(keyringPath)
	if err != nil {
		return []byte(""), err
	}

	key, ok := kr.GetKey(alias)
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

	encoded := []byte(base64.StdEncoding.EncodeToString(encrypted))

	return encoded, nil
}
