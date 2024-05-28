//go:build amd64

package json

import (
	"bytes"
	"github.com/bytedance/sonic"
	"io"
)

func NewDecoder(r io.Reader) sonic.Decoder {
	return sonic.ConfigDefault.NewDecoder(r)
}

func Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func MarshalReader(v interface{}) (io.Reader, error) {
	data, err := sonic.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
