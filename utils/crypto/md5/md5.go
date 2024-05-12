package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func EncodeMD5String(value string) string {
	md5 := md5.Sum([]byte(value))
	return hex.EncodeToString(md5[:])
}

func EncodeMD5(value string) []byte {
	md5 := md5.Sum([]byte(value))
	return md5[:]
}

func MD5String(md5 []byte) string {
	return hex.EncodeToString(md5)
}

func EncodeReaderMD5(r io.Reader) ([]byte, error) {
	hash := md5.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func EncodeReaderMD5String(r io.Reader) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
