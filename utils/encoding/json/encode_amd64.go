//go:build amd64

package json

import (
	"github.com/bytedance/sonic"
)

func Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}
