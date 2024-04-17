//go:build amd64

package json

import (
	"github.com/bytedance/sonic"
	"io"
)

func NewEncoder(w io.Writer) sonic.Encoder {
	return sonic.ConfigDefault.NewEncoder(w)
}

func Unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}
