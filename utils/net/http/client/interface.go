package client

import "net/http"

type MarshalBody interface {
	MarshalBody(contentType string) ([]byte, error)
}

type UnmarshalBody interface {
	UnmarshalBody(contentType string, body []byte) error
}

type MarshalQuery interface {
	MarshalQuery() (string, error)
}

type UnmarshalQuery interface {
	UnmarshalQuery(query string) error
}

type MarshalHeader interface {
	MarshalHeader() ([]string, error)
}

type UnmarshalHeader interface {
	UnmarshalHeader(headers []string) error
}

type SetRequest interface {
	SetRequest(*http.Request)
}

type FromResponse interface {
	FromResponse(response *http.Response)
}
