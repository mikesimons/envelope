package envelope

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"
)

type structuredCodec struct {
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte) (interface{}, error)
}

func resolveInterface(t reflect.Value) reflect.Value {
	if t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr {
		next := t.Elem()
		if !next.IsValid() {
			return reflect.ValueOf(nil)
		}
		return resolveInterface(next)
	}
	return t
}

func codecForFormat(format string) (structuredCodec, error) {
	switch format {
	case "yaml":
		return structuredCodec{
			Marshal: func(v interface{}) ([]byte, error) {
				resolved := resolveInterface(reflect.ValueOf(v))
				if val, ok := (resolved.Interface().(yaml.Node)); ok {
					return yaml.Marshal(&val)
				}

				return yaml.Marshal(v)
			},
			Unmarshal: func(b []byte) (interface{}, error) {
				var out yaml.Node
				err := yaml.Unmarshal(b, &out)
				return out, err
			},
		}, nil
	case "json":
		return structuredCodec{
			Marshal: func(v interface{}) ([]byte, error) {
				return json.MarshalIndent(v, "", "  ")
			},
			Unmarshal: func(b []byte) (interface{}, error) {
				out := make(map[string]interface{})
				err := json.Unmarshal(b, &out)
				return out, err
			},
		}, nil
	case "toml":
		return structuredCodec{
			Marshal: func(v interface{}) ([]byte, error) {
				// Internally we use with map[interface{}]interface{}
				// YAML uses that natively and JSON figures it its actually a map[string]interface{} (or the implementation doesn't care)
				// TOML *does* care so we have to do some type juggling
				ptr := v.(*interface{})
				val := (*ptr).(map[string]interface{})
				return toml.Marshal(val)
			},
			Unmarshal: func(b []byte) (interface{}, error) {
				out := make(map[string]interface{})
				err := toml.Unmarshal(b, &out)
				return out, err
			},
		}, nil
	}

	return structuredCodec{}, fmt.Errorf("Unknown format '%s'", format)
}
