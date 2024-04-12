package json

import (
	"encoding/json"
	"io"
)

func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
