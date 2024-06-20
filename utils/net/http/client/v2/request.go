package client

import (
	"github.com/hopeio/cherry/utils/net/http/client"
)

// Client ...

type Request[RES any] client.Request

func NewRequest[RES any](method, url string) *Request[RES] {
	return &Request[RES]{Method: method, Url: url}
}

func NewFromRequest[RES any](req *client.Request) *Request[RES] {
	return (*Request[RES])(req)
}

func (req *Request[RES]) WithClient(client2 *client.Client) *Request[RES] {
	(*client.Request)(req).WithClient(client2)
	return req
}

func (req *Request[RES]) SetClient(set func(c *client.Client)) *Request[RES] {
	(*client.Request)(req).SetClient(set)
	return req
}

func (req *Request[RES]) Client() *client.Client {
	return (*client.Request)(req).Client()
}
func (r *Request[RES]) Origin() *client.Request {
	return (*client.Request)(r)
}

// Do create a HTTP request
func (r *Request[RES]) Do(param any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Do(param, response)
	return response, err
}
