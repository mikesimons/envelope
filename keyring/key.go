package keyring

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
)

type Key struct {
	Id         uuid.UUID
	Name       string
	Ciphertext []byte
}

// NewKey builds a new key with the given name / encrypted data key
func NewKey(name string, ciphertext []byte) *Key {
	return &Key{
		Id:         uuid.NewV4(),
		Name:       name,
		Ciphertext: ciphertext,
	}
}

func (key *Key) MarshalYAML() (interface{}, error) {
	out := make(map[string]string)
	out["id"] = key.Id.String()
	out["name"] = key.Name
	out["key"] = base64.StdEncoding.EncodeToString(key.Ciphertext)
	return out, nil
}

func (key *Key) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var custom struct {
		Id   uuid.UUID
		Name string
		Key  string
	}

	if err := unmarshal(&custom); err != nil {
		return err
	}

	key.Id = custom.Id
	key.Name = custom.Name
	decoded, err := base64.StdEncoding.DecodeString(string(custom.Key))
	if err != nil {
		return err
	}

	key.Ciphertext = decoded

	return nil
}
