package common

import (
	"encoding/json"
	"errors"
	"github.com/hopeio/cherry/utils/encoding"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"strings"
)

func Unmarshal(format encoding.Format, data []byte, v any) (encoding.Format, error) {
	switch {
	case format == encoding.YAML || format == encoding.YML:
		return format, yaml.Unmarshal(data, v)
	case format == encoding.TOML:
		return format, toml.Unmarshal(data, v)
	case format == encoding.JSON:
		return format, json.Unmarshal(data, v)
	default:
		if err := toml.Unmarshal(data, v); err == nil {
			return encoding.TOML, nil
		}
		if err := json.Unmarshal(data, v); err == nil {
			return encoding.JSON, nil
		} else if strings.Contains(err.Error(), "json: unknown field") {
			return encoding.JSON, err
		}

		yamlError := yaml.Unmarshal(data, v)

		if yamlError == nil {
			return encoding.YAML, nil
		} else if yErr, ok := yamlError.(*yaml.TypeError); ok {
			return encoding.YAML, yErr
		}

		return "", errors.New("failed to decode")
	}
}

func Marshal(format encoding.Format, v any) ([]byte, error) {
	switch {
	case format == encoding.YAML || format == encoding.YML:
		return yaml.Marshal(v)
	case format == encoding.TOML:
		return toml.Marshal(v)
	case format == encoding.JSON:
		return json.Marshal(v)
	default:
		return toml.Marshal(v)
	}
}
