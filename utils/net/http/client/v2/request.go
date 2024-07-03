package client

import (
	"context"
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

func (req *Request[RES]) AddHeader(k, v string) *Request[RES] {
	(*client.Request)(req).AddHeader(k, v)
	return req
}

func (req *Request[RES]) ContentType(contentType client.ContentType) *Request[RES] {
	(*client.Request)(req).ContentType(contentType)
	return req
}

func (req *Request[RES]) Context(ctx context.Context) *Request[RES] {
	(*client.Request)(req).Context(ctx)
	return req
}

func (req *Request[RES]) DoNoParam() (*RES, error) {
	response := new(RES)
	return response, (*client.Request)(req).Do(nil, response)
}

// Do create a HTTP request
func (r *Request[RES]) Do(param any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Do(param, response)
	return response, err
}
