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
	req.Client = client2
	return req
}

func (r *Request[RES]) Origin() *client.Request {
	return (*client.Request)(r)
}

// Do create a HTTP request
func (r *Request[RES]) Do(param any) (*RES, error) {
	response := new(RES)
	if r.Client == nil {
		r.Client = &client.Client{}
	}
	err := r.Client.Do(r.Method, r.Url, param, response)
	return response, err
}

func (req *Request[RES]) Clone() *Request[RES] {
	newReq := &(*req)
	newReq.Client = req.Client.Clone()
	return newReq
}
