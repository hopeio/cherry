//go:build amd64 && !sonic

package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"io"
)

var Standard = jsoniter.ConfigCompatibleWithStandardLibrary

func SupportPrivateFields() {
	extra.SupportPrivateFields()
}

var WithPrivateField = jsoniter.Config{
	IndentionStep:                 4,
	MarshalFloatWith6Digits:       true,
	EscapeHTML:                    true,
	SortMapKeys:                   true,
	UseNumber:                     true,
	ObjectFieldMustBeSimpleString: true,
}.Froze()

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return Standard.NewEncoder(w)
}

func Marshal(v interface{}) ([]byte, error) {
	return Standard.Marshal(v)
}

func MarshalString(v any) (string, error) {
	return Standard.MarshalToString(v)
}

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return Standard.NewDecoder(r)
}

func Unmarshal(data []byte, v any) error {
	return jsoniter.ConfigDefault.Unmarshal(data, v)
}

func UnmarshalString(data string, v any) error {
	return jsoniter.ConfigDefault.UnmarshalFromString(data, v)
}
