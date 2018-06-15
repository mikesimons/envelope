package keyring

import (
	"io/ioutil"
	"os"

	"github.com/ansel1/merry"
	"github.com/mikesimons/envelope/keysvc"
	"github.com/satori/go.uuid"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type YAMLKeyring struct {
	File string
	Keys []*keysvc.Key
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
		return keyring, merry.Wrap(err).WithUserMessage("Couldn't open keyring file")
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return keyring, merry.Wrap(err).WithUserMessage("Couldn't read keyring file")
	}

	err = yaml.Unmarshal(bytes, &keyring.Keys)

	if err != nil {
		return keyring, merry.Wrap(err).WithUserMessage("Couldn't parse YAML keyring data")
	}

	return keyring, nil
}

func (kr *YAMLKeyring) Decrypt(data []byte) ([]byte, error) {
	decoded, err := keysvc.DecodeEncrypted(data)
	if err != nil {
		return []byte(""), err
	}

	key, ok := kr.GetKey(decoded.KeyID.String())
	if !ok {
		return []byte(""), merry.New("Couldn't find key").WithValue("keyring id", decoded.KeyID.String())
	}

	return key.DecryptEncryptedItem(decoded)
}

// GetKeys returns an array of all keys in the keyring
func (kr *YAMLKeyring) GetKeys() []*keysvc.Key {
	return kr.Keys
}

// AddKey adds a predefined key to the keyring
func (kr *YAMLKeyring) AddKey(key *keysvc.Key) error {
	_, idExists := kr.GetKey(key.ID.String())
	_, aliasExists := kr.GetKey(key.Alias)

	if idExists || aliasExists {
		return merry.New("Couldn't add key because it clashes with an existing key alias or id").
			WithValue("keyring id", key.ID.String()).
			WithValue("keyring alias", key.Alias)
	}

	kr.Keys = append(kr.Keys, key)

	updated, err := yaml.Marshal(kr.Keys)
	if err != nil {
		return merry.Wrap(err).WithUserMessage("Couldn't add key")
	}

	err = afero.WriteFile(Fs, kr.File, updated, 0644)

	if err != nil {
		return merry.Wrap(err).WithUserMessage("Couldn't add key")
	}
	return nil
}

// GetKey gets an individual key from the keyring
func (kr *YAMLKeyring) GetKey(aliasOrID string) (*keysvc.Key, bool) {
	id := uuid.FromStringOrNil(aliasOrID)
	for _, k := range kr.Keys {
		if (id != uuid.UUID{} && k.ID == id) || k.Alias == aliasOrID {
			return k, true
		}
	}
	return nil, false
}
