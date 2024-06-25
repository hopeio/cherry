//go:build amd64 && snoic

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

func MarshalString(v any) (string, error) {
	return sonic.MarshalString(v)
}

func NewEncoder(w io.Writer) sonic.Encoder {
	return sonic.ConfigDefault.NewEncoder(w)
}

func Unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

func UnmarshalString(data string, v any) error {
	return sonic.UnmarshalString(data, v)
}
