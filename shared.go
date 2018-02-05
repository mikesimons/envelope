package sekrits

import (
	"encoding/json"
	"fmt"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v2"
)

type structuredCodec struct {
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

func codecForFormat(format string) (structuredCodec, error) {
	switch format {
	case "yaml":
		return structuredCodec{
			Marshal:   yaml.Marshal,
			Unmarshal: yaml.Unmarshal,
		}, nil
	case "json":
		return structuredCodec{
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		}, nil
	case "toml":
		return structuredCodec{
			Marshal:   toml.Marshal,
			Unmarshal: toml.Unmarshal,
		}, nil
	}

	return structuredCodec{}, fmt.Errorf("Unknown format '%s'", format)
}