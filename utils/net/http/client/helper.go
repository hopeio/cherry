package client

import (
	"io"
	"net/http"
)

func DefaultHeaderRequest() *Request {
	req := newRequest("", "")
	req.Header(DefaultHeader())
	return req
}

func NewGet(url string) *Request {
	return newRequest(http.MethodGet, url)
}

func NewPost(url string) *Request {
	return newRequest(http.MethodPost, url)
}

func NewPut(url string) *Request {
	return newRequest(http.MethodPut, url)
}

func NewDelete(url string) *Request {
	return newRequest(http.MethodDelete, url)
}

func Get(url string, response any) error {
	return NewGet(url).DisableLog().DoNoParam(response)
}

func GetStream(url string) (io.ReadCloser, error) {
	var resp *http.Response
	err := Get(url, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func Post(url string, param, response interface{}) error {
	return NewPost(url).DisableLog().Do(param, response)
}

func Put(url string, param, response interface{}) error {
	return NewPut(url).DisableLog().Do(param, response)
}

func Delete(url string, param, response interface{}) error {
	return NewDelete(url).DisableLog().Do(param, response)
}
