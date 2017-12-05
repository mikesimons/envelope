package keyring

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
)

type Key struct {
	Id         uuid.UUID
	Alias      string
	Ciphertext []byte
}

// NewKey builds a new key with the given alias / encrypted data key
func NewKey(alias string, ciphertext []byte) *Key {
	return &Key{
		Id:         uuid.NewV4(),
		Alias:      alias,
		Ciphertext: ciphertext,
	}
}

func (key *Key) MarshalYAML() (interface{}, error) {
	out := make(map[string]string)
	out["id"] = key.Id.String()
	out["alias"] = key.Alias
	out["key"] = base64.StdEncoding.EncodeToString(key.Ciphertext)
	return out, nil
}

func (key *Key) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var custom struct {
		Id    uuid.UUID
		Alias string
		Key   string
	}

	if err := unmarshal(&custom); err != nil {
		return err
	}

	key.Id = custom.Id
	key.Alias = custom.Alias
	decoded, err := base64.StdEncoding.DecodeString(string(custom.Key))
	if err != nil {
		return err
	}

	key.Ciphertext = decoded

	return nil
}
