package sekrits

import (
	"github.com/mikesimons/sekrits/datakey/provider"
	"github.com/mikesimons/sekrits/keyring"
	"net/url"
)

func AddKey(keyringPath string, alias string, providerDsn string) (string, error) {
	kr, err := keyring.Load(keyringPath)
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(providerDsn)
	if err != nil {
		return "", err
	}

	dkp, err := provider.Factory(parsed.Scheme)
	if err != nil {
		return "", err
	}

	rawKey, err := dkp.GenerateDatakey(parsed.Host)
	if err != nil {
		return "", err
	}

	key := keyring.NewKey(alias, rawKey)
	kr.AddKey(key)

	return key.Id.String(), nil
}
