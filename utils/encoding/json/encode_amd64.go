//go:build amd64

package json

import (
	"github.com/bytedance/sonic"
	"io"
)

func NewDecoder(r io.Reader) sonic.Decoder {
	return sonic.ConfigDefault.NewDecoder(r)
}

func Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}
