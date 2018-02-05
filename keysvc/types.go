package keysvc

import (
	"github.com/satori/go.uuid"
)

type encryptedData struct {
	KeyId      uuid.UUID
	Ciphertext []byte
}

type KeyServiceProvider interface {
	GenerateDatakey(master string, context map[string]string) ([]byte, error)
	DecryptDatakey(ciphertext *[]byte, plaintext *[]byte, context map[string]string) error
}

type Key struct {
	Id         uuid.UUID
	Alias      string
	Ciphertext []byte
	Plaintext  []byte
	Type       string
	Context    map[string]string
}
