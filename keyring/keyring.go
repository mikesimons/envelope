package keyring

import (
	"github.com/mikesimons/envelope/keysvc"
	"github.com/spf13/afero"
)

var Fs = afero.NewOsFs()

type Keyring interface {
	GetKeys() []*keysvc.Key
	GetKey(aliasOrID string) (*keysvc.Key, bool)
	AddKey(*keysvc.Key) error
	Decrypt([]byte) ([]byte, error)
}

// Load loads keyring data (only YAML supported right now)
func Load(keyring string) (Keyring, error) {
	return LoadYAML(keyring)
}
