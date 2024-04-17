//go:build !amd64

package json

import (
	jsoniter "github.com/json-iterator/go"
	"io"
)

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return jsoniter.ConfigDefault.NewEncoder(w)
}

func Unmarshal(data []byte, v any) error {
	return jsoniter.ConfigDefault.Unmarshal(data, v)
}
