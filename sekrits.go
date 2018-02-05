package sekrits

import (
	"github.com/mikesimons/sekrits/keyring"
)

type Sekrits struct {
	Keyring keyring.Keyring
	Prefix  string
}

func WithYamlKeyring(path string) (*Sekrits, error) {
	kr, err := keyring.Load(path)
	if err != nil {
		return &Sekrits{}, err
	}

	return &Sekrits{
		Keyring: kr,
		Prefix:  "!!sekrit:",
	}, nil
}
