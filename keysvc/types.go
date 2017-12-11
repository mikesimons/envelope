package keysvc

import (
	"github.com/satori/go.uuid"
)

type encryptedData struct {
	KeyId      uuid.UUID
	Ciphertext []byte
}

type KeyServiceProvider interface {
	GenerateDatakey(master string) ([]byte, error)
	DecryptDatakey(ciphertext *[]byte, plaintext *[]byte) error
}

type Key struct {
	Id         uuid.UUID
	Alias      string
	Ciphertext []byte
	Plaintext  []byte
	Type       string
}
