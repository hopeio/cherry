package client

import "github.com/hopeio/cherry/utils/net/http/client"

type ResponseInterface[T any] interface {
	client.ResponseBodyCheck
	GetData() T
}

// 一个语法糖，一般不用
type SubDataRequest[RES ResponseInterface[T], T any] Request[RES]

func NewSubDataRequestParams[RES ResponseInterface[T], T any](req *client.Request) *SubDataRequest[RES, T] {
	return (*SubDataRequest[RES, T])(req)
}

func (req *SubDataRequest[RES, T]) Origin() *client.Request {
	return (*client.Request)(req)
}

// Do create a HTTP request
func (r *SubDataRequest[RES, T]) Do(req any) (T, error) {
	var response RES
	err := (*client.Request)(r).Do(req, response)
	if err != nil {
		return response.GetData(), err
	}
	return response.GetData(), err
}

func (req *SubDataRequest[RES, T]) Get(url string) (T, error) {
	var response RES
	err := (*client.Request)(req).Url(url).Do(req, &response)
	if err != nil {
		return response.GetData(), err
	}

	return response.GetData(), nil
}
