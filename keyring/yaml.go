package keyring

import (
	"fmt"
	errors "github.com/hashicorp/errwrap"
	"github.com/satori/go.uuid"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type YAMLKeyring struct {
	File string
	Keys []*Key
}

// LoadYAML loads a YAML file as a keyring
// If the file does not exist it will return an empty keyring but not an error.
// This is a convenience function to avoid having a separate "create keyring" step for new users
func LoadYAML(path string) (Keyring, error) {
	keyring := &YAMLKeyring{File: path}

	_, err := Fs.Stat(path)

	// If the file does not exist we assume the user wants a new keyring file.
	// If this is not what the user wants, no other action will do anything useful
	// as they all need keys that are not loaded to function (except add-key)
	if os.IsNotExist(err) {
		return keyring, nil
	}

	file, err := Fs.Open(path)
	if err != nil {
		return keyring, errors.Wrapf("Couldn't open keyring file", err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return keyring, errors.Wrapf("Couldn't read keyring file", err)
	}

	err = yaml.Unmarshal(bytes, &keyring.Keys)

	if err != nil {
		return keyring, errors.Wrapf("Couldn't parse YAML keyring data", err)
	}

	return keyring, nil
}

// GetKeys returns an array of all keys in the keyring
func (kr *YAMLKeyring) GetKeys() []*Key {
	return kr.Keys
}

// AddKey adds a predefined key to the keyring
func (kr *YAMLKeyring) AddKey(key *Key) error {
	_, idExists := kr.GetKey(key.Id.String())
	_, nameExists := kr.GetKey(key.Name)

	if idExists || nameExists {
		return fmt.Errorf("Couldn't add key because '%s' clashes with an existing key alias or id", key.Name)
	}

	kr.Keys = append(kr.Keys, key)

	updated, err := yaml.Marshal(kr.Keys)
	if err != nil {
		return errors.Wrapf("Couldn't add key", err)
	}

	err = afero.WriteFile(Fs, kr.File, updated, 0644)

	if err != nil {
		return errors.Wrapf("Couldn't add key", err)
	}
	return nil
}

// GetKey gets an individual key from the keyring
func (kr *YAMLKeyring) GetKey(nameOrId string) (*Key, bool) {
	id := uuid.FromStringOrNil(nameOrId)
	for _, k := range kr.Keys {
		if (id != uuid.UUID{} && k.Id == id) || k.Name == nameOrId {
			return k, true
		}
	}
	return &Key{}, false
}
