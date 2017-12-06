package datakey

import (
	"fmt"
)

type DatakeyProvider interface {
	GenerateDatakey(string) ([]byte, error)
	DecryptDatakey(key []byte)
}

func Factory(name string) (DatakeyProvider, error) {
	if name == "awskms" {
		return NewAWSKMS()
	}

	return nil, fmt.Errorf("Unknown datakey provider: %s", name)
}
