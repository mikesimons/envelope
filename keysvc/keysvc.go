package keysvc

import (
	"fmt"
	errors "github.com/hashicorp/errwrap"
	"github.com/mikesimons/sekrits/keysvc/awskms"
	"gopkg.in/mgo.v2/bson"
	"net/url"
)

func awsKeySvcFn() (KeyServiceProvider, error) {
	return awskms.New()
}

var services = map[string]func() (KeyServiceProvider, error){
	"awskms": awsKeySvcFn,
}

func AddKeyServiceFn(name string, fn func() (KeyServiceProvider, error)) {
	services[name] = fn
}

func GetKeyService(name string) (KeyServiceProvider, error) {
	if fn, ok := services[name]; ok {
		return fn()
	}

	return nil, fmt.Errorf("Unknown key service: %s", name)
}

func GenerateDatakey(alias string, masterKey string, context map[string]string) (*Key, error) {
	parsed, err := url.Parse(masterKey)
	if err != nil {
		return nil, errors.Wrapf("Could not parse master key URL", err)
	}

	keysvc, err := GetKeyService(parsed.Scheme)
	if err != nil {
		return nil, errors.Wrapf("Could not initialize key service", err)
	}

	ciphertext, err := keysvc.GenerateDatakey(fmt.Sprintf("%s%s", parsed.Host, parsed.Path), context)
	if err != nil {
		return nil, errors.Wrapf("Could not generate data key", err)
	}

	return NewKey(alias, parsed.Scheme, ciphertext, context), nil
}

func DecodeEncrypted(data []byte) (encryptedData, error) {
	var encrypted encryptedData
	err := bson.Unmarshal(data, &encrypted)
	if err != nil {
		return encryptedData{}, errors.Wrapf("Could not decode an encrypted item; it is possibly corrupted", err)
	}
	return encrypted, nil
}
