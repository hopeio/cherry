//go:build amd64

package json

import (
	"github.com/bytedance/sonic"
)

func Unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}
