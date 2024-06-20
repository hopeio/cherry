package http

import (
	"net/http"
)

type Header []string

func NewHeader() *Header {
	h := make(Header, 0, 6)
	return &h
}

func (h *Header) Add(k, v string) *Header {
	*h = append(*h, k, v)
	return h
}

func (h *Header) Set(k, v string) *Header {
	header := *h
	for i, s := range header {
		if s == k {
			header[i+1] = v
			return h
		}
	}
	return h.Add(k, v)
}

func (h *Header) IntoHttpHeader(header http.Header) {
	res := *h
	hlen := len(res)
	for i := 0; i < hlen && i+1 < hlen; i += 2 {
		header.Set(res[i], res[i+1])
	}
}

func (h Header) Clone() Header {
	newh := make(Header, len(h))
	copy(newh, h)
	return newh
}
