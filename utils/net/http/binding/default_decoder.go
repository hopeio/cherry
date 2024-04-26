package binding

import (
	reflecti "github.com/hopeio/cherry/utils/reflect/converter"
	"reflect"
)

var defaultDecoder = NewDecoder()

func DefaultDecoder() *Decoder {
	return defaultDecoder
}

func SetAliasTag(tag string) {
	defaultDecoder.SetAliasTag(tag)
}

func ZeroEmpty(z bool) {
	defaultDecoder.zeroEmpty = z
}

func IgnoreUnknownKeys(i bool) {
	defaultDecoder.ignoreUnknownKeys = i
}

// RegisterConverter registers a converter function for a custom type.
func RegisterConverter(value interface{}, converterFunc reflecti.StringConverter) {
	defaultDecoder.cache.registerConverter(value, converterFunc)
}

func Decode(dst interface{}, src map[string][]string) error {
	return defaultDecoder.Decode(dst, src)
}

func PickDecode(v reflect.Value, src map[string][]string) error {
	return defaultDecoder.PickDecode(v, src)
}
