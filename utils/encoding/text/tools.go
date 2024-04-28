package text

import (
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

func DetermineEncoding(content []byte, contentType string) (e encoding.Encoding, name string, certain bool) {
	return charset.DetermineEncoding(content, contentType)
}
