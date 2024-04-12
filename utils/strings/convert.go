package stringsi

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"strings"
	"unsafe"
)

func ToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// GBK 转 UTF-8
func GBKToUTF8(s string) (string, error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return ToString(b), nil
}

func BytesGBKToUTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}

// UTF-8 转 GBK

func UTF8ToGBK(s string) (string, error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return ToString(b), nil
}

func BytesUTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	return io.ReadAll(reader)
}
