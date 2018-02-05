package sekrits

import (
	"github.com/mikesimons/sekrits/keysvc"
)

func (s *Sekrits) AddKey(alias string, masterKey string, context map[string]string) (string, error) {
	key, err := keysvc.GenerateDatakey(alias, masterKey, context)
	if err != nil {
		return "", err
	}

	err = s.Keyring.AddKey(key)
	if err != nil {
		return "", err
	}

	return key.Id.String(), nil
}
