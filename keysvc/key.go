package keysvc

import (
	"encoding/base64"
	"fmt"
	//errors "github.com/hashicorp/errwrap"
	"math/rand"

	"github.com/ansel1/merry"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/nacl/secretbox"
	"gopkg.in/mgo.v2/bson"
)

// NewKey builds a new key with the given alias / encrypted data key
func NewKey(alias string, keyType string, ciphertext []byte, context map[string]string) *Key {
	return &Key{
		Id:         uuid.NewV4(),
		Alias:      alias,
		Ciphertext: ciphertext,
		Type:       keyType,
		Context:    context,
	}
}

func (key *Key) decryptDatakey() error {
	keysvc, err := GetKeyService(key.Type)
	if err != nil {
		return err
	}

	err = keysvc.DecryptDatakey(&key.Ciphertext, &key.Plaintext, key.Context)
	if err != nil {
		return merry.Wrap(err).
			WithValue("keyring id", key.Id.String()).
			WithValue("keyring alias", key.Alias)
	}

	return nil
}

func (key *Key) EncryptBytes(data []byte) ([]byte, error) {
	err := key.decryptDatakey()
	if err != nil {
		return []byte(""), err
	}

	var nonce [24]byte
	rand.Read(nonce[:])

	ret := &encryptedData{
		KeyId: key.Id,
	}

	var plaintext [32]byte
	copy(plaintext[:], key.Plaintext[:32])
	ret.Ciphertext = secretbox.Seal(nonce[:], data, &nonce, &plaintext)
	return bson.Marshal(ret)
}

func (key *Key) Decrypt(data []byte) ([]byte, error) {
	encrypted, err := DecodeEncrypted(data)
	if err != nil {
		return []byte(""), err
	}
	return key.DecryptEncryptedItem(encrypted)
}

func (key *Key) DecryptEncryptedItem(encrypted encryptedData) ([]byte, error) {
	err := key.decryptDatakey()
	if err != nil {
		return []byte(""), err
	}

	var nonce [24]byte
	var plaintext [32]byte
	copy(nonce[:], encrypted.Ciphertext[:24])
	copy(plaintext[:], key.Plaintext[:32])
	decrypted, ok := secretbox.Open(nil, encrypted.Ciphertext[24:], &nonce, &plaintext)

	if !ok {
		return []byte(""), fmt.Errorf("Could not decrypt secret with data key; the secret may be corrupted")
	}

	return decrypted, nil
}

func (key *Key) MarshalYAML() (interface{}, error) {
	out := make(map[string]interface{})
	out["id"] = key.Id.String()
	out["alias"] = key.Alias
	out["key"] = base64.StdEncoding.EncodeToString(key.Ciphertext)
	out["type"] = key.Type
	out["context"] = key.Context
	return out, nil
}

func (key *Key) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var custom struct {
		Id      uuid.UUID
		Alias   string
		Key     string
		Type    string
		Context map[string]string
	}

	if err := unmarshal(&custom); err != nil {
		return err
	}

	key.Id = custom.Id
	key.Alias = custom.Alias
	key.Type = custom.Type
	key.Context = custom.Context
	decoded, err := base64.StdEncoding.DecodeString(string(custom.Key))
	if err != nil {
		return err
	}

	key.Ciphertext = decoded

	return nil
}
