package keyring

import (
	"github.com/spf13/afero"
)

var Fs = afero.NewOsFs()

type Keyring interface {
	GetKeys() []*Key
	GetKey(aliasOrId string) (*Key, bool)
	AddKey(*Key) error
}

// Load loads keyring data (only YAML supported right now)
func Load(keyring string) (Keyring, error) {
	return LoadYAML(keyring)
}
