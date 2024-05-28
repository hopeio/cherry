package client

import (
	"io"
	"net/http"
)

func NewGetRequest(url string) *Request {
	return newRequest(url, http.MethodGet)
}

func NewPostRequest(url string) *Request {
	return newRequest(url, http.MethodPost)
}

func NewPutRequest(url string) *Request {
	return newRequest(url, http.MethodPut)
}

func NewDeleteRequest(url string) *Request {
	return newRequest(url, http.MethodDelete)
}

func DoGet(url string, response any) error {
	return NewGetRequest(url).DisableLog().DoWithNoParam(response)
}

func DoGetStream(url string) (io.ReadCloser, error) {
	var resp *http.Response
	err := DoGet(url, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func DoPost(url string, param, response interface{}) error {
	return NewPostRequest(url).DisableLog().Do(param, response)
}

func DoPut(url string, param, response interface{}) error {
	return NewPutRequest(url).DisableLog().Do(param, response)
}

func DoDelete(url string, param, response interface{}) error {
	return NewDeleteRequest(url).DisableLog().Do(param, response)
}
