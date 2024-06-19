package client

import (
	"io"
	"net/http"
)

type Request struct {
	Method, Url string
	*Client
}

func NewRequest(method, url string) *Request {
	return &Request{
		Method: method,
		Url:    url,
		Client: New(),
	}
}

func (req *Request) Do(param, response interface{}) error {
	return req.Client.Do(req.Method, req.Url, param, response)
}

func (req *Request) DoEmpty() error {
	return req.Client.Do(req.Method, req.Url, nil, nil)
}

func (req *Request) DoNoParam(response interface{}) error {
	return req.Client.Do(req.Method, req.Url, nil, response)
}

func (req *Request) DoNoResponse(param interface{}) error {
	return req.Client.Do(req.Method, req.Url, param, nil)
}

func (req *Request) DoRaw(param interface{}) (RawBytes, error) {
	var raw RawBytes
	err := req.Client.Do(req.Method, req.Url, param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func (req *Request) DoStream(method, url string, param interface{}) (io.ReadCloser, error) {
	var resp *http.Response
	err := req.Client.Do(method, url, param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (req *Request) Get(param, response interface{}) error {
	return req.Client.Do(http.MethodGet, req.Url, param, response)
}

func (req *Request) Post(param, response interface{}) error {
	return req.Client.Do(http.MethodPost, req.Url, param, response)
}

func (req *Request) Put(param, response interface{}) error {
	return req.Client.Do(http.MethodPut, req.Url, param, response)
}

func (req *Request) Delete(param, response interface{}) error {
	return req.Client.Do(http.MethodDelete, req.Url, param, response)
}
