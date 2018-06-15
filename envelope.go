package envelope

import (
	"github.com/mikesimons/envelope/keyring"
	"github.com/mikesimons/traverser"
)

// Envelope holds the configuration and keyring
type Envelope struct {
	Keyring                  keyring.Keyring
	Prefix                   string
	StructuredErrorBehaviour func(error) (traverser.Op, error)
}

// WithYamlKeyring initializes Envelope with a yaml file based keyring
func WithYamlKeyring(path string) (*Envelope, error) {
	kr, err := keyring.Load(path)
	if err != nil {
		return &Envelope{}, err
	}

	return &Envelope{
		Keyring: kr,
		Prefix:  "!!enveloped:",
		StructuredErrorBehaviour: func(e error) (traverser.Op, error) { return traverser.Noop() },
	}, nil
}
