package client

import (
	"io"
	"net/http"
)

func DefaultHeaderClient() *Client {
	return newClient().Header(DefaultHeader())
}

func DefaultHeaderRequest() *Request {
	return &Request{Client: newClient().Header(DefaultHeader())}
}

func NewGet(url string) *Request {
	return NewRequest(http.MethodGet, url)
}

func NewPost(url string) *Request {
	return NewRequest(http.MethodPost, url)
}

func NewPut(url string) *Request {
	return NewRequest(http.MethodPut, url)
}

func NewDelete(url string) *Request {
	return NewRequest(http.MethodDelete, url)
}

func Get(url string, param, response any) error {
	return NewGet(url).Do(param, response)
}

func GetNP(url string, response any) error {
	return Get(url, nil, response)
}

func GetStream(url string, param any) (io.ReadCloser, error) {
	var resp *http.Response
	err := Get(url, param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func GetStreamNP(url string) (io.ReadCloser, error) {
	return GetStream(url, nil)
}

func Post(url string, param, response interface{}) error {
	return New().DisableLog().Post(url, param, response)
}

func Put(url string, param, response interface{}) error {
	return New().DisableLog().Put(url, param, response)
}

func Delete(url string, param, response interface{}) error {
	return New().DisableLog().Delete(url, param, response)
}
