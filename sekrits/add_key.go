package sekrits

import (
	"fmt"
	"github.com/mikesimons/sekrits/keyring"
	"github.com/mikesimons/sekrits/keysvc"
	"net/url"
)

func AddKey(keyringPath string, alias string, datakeyDsn string) (string, error) {
	kr, err := keyring.Load(keyringPath)
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(datakeyDsn)
	if err != nil {
		return "", err
	}

	dkp, err := keysvc.Factory(parsed.Scheme)
	if err != nil {
		return "", err
	}

	rawKey, err := dkp.GenerateDatakey(fmt.Sprintf("%s%s", parsed.Host, parsed.Path))
	if err != nil {
		return "", err
	}

	key := keyring.NewKey(alias, rawKey)
	err = kr.AddKey(key)
	if err != nil {
		return "", err
	}

	return key.Id.String(), nil
}
