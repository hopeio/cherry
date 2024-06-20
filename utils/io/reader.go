package io

import (
	"io"
)

func ReadReadCloser(readCloser io.ReadCloser) ([]byte, error) {
	data, err := io.ReadAll(readCloser)
	if err != nil {
		return nil, err
	}
	readCloser.Close()
	return data, nil
}

type warpCloser struct {
	io.Reader
}

func (*warpCloser) Close() error {
	return nil
}

func WarpCloser(body io.Reader) io.ReadCloser {
	return &warpCloser{body}
}
