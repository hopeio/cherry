//go:build !amd64

package json

func Unmarshal(data []byte, v any) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, v)
}
