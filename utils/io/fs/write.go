package fs

import (
	"bytes"
)

func WriteBuffer(buf *bytes.Buffer, filename string) (n int, err error) {
	f, _ := Create(filename)
	defer f.Close()
	return f.Write(buf.Bytes())
}

func Write(data []byte, filename string) (n int, err error) {
	f, _ := Create(filename)
	defer f.Close()
	return f.Write(data)
}
