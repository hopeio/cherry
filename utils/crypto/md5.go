package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5String(value string) string {
	md5 := md5.Sum([]byte(value))
	return hex.EncodeToString(md5[:])
}

func EncodeMD5Bytes(value string) []byte {
	md5 := md5.Sum([]byte(value))
	return md5[:]
}
