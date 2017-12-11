package sekrits

import (
	"github.com/mikesimons/sekrits/keyring"
	"github.com/mikesimons/sekrits/keysvc"
)

func AddKey(keyringPath string, alias string, masterKey string) (string, error) {
	kr, err := keyring.Load(keyringPath)
	if err != nil {
		return "", err
	}

	key, err := keysvc.GenerateDatakey(alias, masterKey)
	if err != nil {
		return "", err
	}

	err = kr.AddKey(key)
	if err != nil {
		return "", err
	}

	return key.Id.String(), nil
}
