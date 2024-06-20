package client

import (
	"io"
	"net/http"
)

func DefaultHeaderClient() *Client {
	return New().Header(DefaultHeader())
}

func DefaultHeaderRequest() *Request {
	return &Request{client: New().Header(DefaultHeader())}
}

func GetRequest(url string) *Request {
	return NewRequest(http.MethodGet, url)
}

func PostRequest(url string) *Request {
	return NewRequest(http.MethodPost, url)
}

func PutRequest(url string) *Request {
	return NewRequest(http.MethodPut, url)
}

func DeleteRequest(url string) *Request {
	return NewRequest(http.MethodDelete, url)
}

func Get(url string, param, response any) error {
	return GetRequest(url).Do(param, response)
}

func GetX(url string, response any) error {
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

func GetStreamX(url string) (io.ReadCloser, error) {
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
