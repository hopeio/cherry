package client

import (
	"github.com/hopeio/cherry/utils/net/http/client"
)

type Request[RES any] client.Request

func NewRequest[RES any](method, url string) *Request[RES] {
	return (*Request[RES])(client.NewRequest(method, url))
}

func (req *Request[RES]) Do(param any) error {
	response := new(RES)
	return (*client.Request)(req).Do(param, response)
}

func (req *Request[RES]) Get(param any) error {
	response := new(RES)
	return (*client.Request)(req).Get(param, response)
}

func (req *Request[RES]) Post(param any) error {
	response := new(RES)
	return (*client.Request)(req).Post(param, response)
}

func (req *Request[RES]) Put(param any) error {
	response := new(RES)
	return (*client.Request)(req).Put(param, response)
}

func (req *Request[RES]) Delete(param any) error {
	response := new(RES)
	return (*client.Request)(req).Delete(param, response)
}
