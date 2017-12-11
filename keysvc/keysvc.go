package keysvc

import (
	"fmt"
	"github.com/mikesimons/sekrits/keysvc/awskms"
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

func GenerateDatakey(alias string, masterKey string) (*Key, error) {
	parsed, err := url.Parse(masterKey)
	if err != nil {
		return nil, err
	}

	keysvc, err := GetKeyService(parsed.Scheme)
	if err != nil {
		return nil, err
	}

	ciphertext, err := keysvc.GenerateDatakey(fmt.Sprintf("%s%s", parsed.Host, parsed.Path))
	if err != nil {
		return nil, err
	}

	return NewKey(alias, parsed.Scheme, ciphertext), nil
}
