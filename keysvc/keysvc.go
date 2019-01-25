package keysvc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ansel1/merry"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mikesimons/envelope/keysvc/awskms"
	"gopkg.in/mgo.v2/bson"
)

func awsKeySvcFn() (KeyServiceProvider, error) {
	awsSession := session.Must(session.NewSession())
	return awskms.New(awsSession)
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

	var known []string
	for k := range services {
		known = append(known, k)
	}
	return nil, merry.Errorf("Unknown key service. Must be one of: %s", strings.Join(known, ", ")).WithValue("svc", name)
}

func GenerateDatakey(alias string, masterKey string, context map[string]string) (*Key, error) {
	parsed, err := url.Parse(masterKey)
	if err != nil {
		return nil, merry.Wrap(err).
			WithUserMessage("Could not parse master key URL").
			WithValue("master key", masterKey).
			WithValue("alias", alias)
	}

	keysvc, err := GetKeyService(parsed.Scheme)
	if err != nil {
		return nil, merry.Wrap(err).
			WithUserMessage("Could not initialize key service").
			WithValue("master key", masterKey).
			WithValue("alias", alias)
	}

	ciphertext, err := keysvc.GenerateDatakey(fmt.Sprintf("%s%s", parsed.Host, parsed.Path), context)
	if err != nil {
		return nil, merry.Wrap(err).
			WithValue("master key", masterKey).
			WithValue("alias", alias)
	}

	return NewKey(alias, parsed.Scheme, ciphertext, context), nil
}

func DecodeEncrypted(data []byte) (encryptedData, error) {
	var encrypted encryptedData
	err := bson.Unmarshal(data, &encrypted)
	if err != nil {
		return encryptedData{}, merry.Wrap(err).
			WithUserMessage("Could not decode an encrypted item; it is possibly corrupted")
	}
	return encrypted, nil
}
