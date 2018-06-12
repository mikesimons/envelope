package sekrits

import (
	"github.com/mikesimons/sekrits/keyring"
	"github.com/mikesimons/traverser"
)

type Sekrits struct {
	Keyring                  keyring.Keyring
	Prefix                   string
	StructuredErrorBehaviour func(error) (traverser.Op, error)
}

func WithYamlKeyring(path string) (*Sekrits, error) {
	kr, err := keyring.Load(path)
	if err != nil {
		return &Sekrits{}, err
	}

	return &Sekrits{
		Keyring: kr,
		Prefix:  "!!sekrit:",
		StructuredErrorBehaviour: func(e error) (traverser.Op, error) { return traverser.Noop() },
	}, nil
}
