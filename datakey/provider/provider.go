package provider

import (
	"fmt"
	"github.com/mikesimons/sekrits/datakey/provider/awskms"
)

type DatakeyProvider interface {
	GenerateDatakey(string) ([]byte, error)
	DecryptDatakey(key []byte)
}

func Factory(name string) (DatakeyProvider, error) {
	if name == "awskms" {
		return awskms.New()
	}

	return nil, fmt.Errorf("Unknown datakey provider: %s", name)
}
