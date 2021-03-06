package envelope

import (
	"github.com/mikesimons/envelope/keysvc"
)

// AddKey will add the given key to the keyring with alias & context
func (s *Envelope) AddKey(alias string, masterKey string, context map[string]string) (string, error) {
	key, err := keysvc.GenerateDatakey(alias, masterKey, context)
	if err != nil {
		return "", err
	}

	err = s.Keyring.AddKey(key)
	if err != nil {
		return "", err
	}

	return key.ID.String(), nil
}
