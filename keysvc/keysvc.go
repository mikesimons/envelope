package keysvc

import (
	"fmt"
)

type KeyServiceProvider interface {
	GenerateDatakey(string) ([]byte, error)
	DecryptDatakey(key []byte)
}

func Factory(name string) (KeyServiceProvider, error) {
	if name == "awskms" {
		return NewAWSKMS()
	}

	return nil, fmt.Errorf("Unknown key service provider: %s", name)
}
