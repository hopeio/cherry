package client

import (
	"github.com/hopeio/cherry/utils/net/http/client"
	"net/http"
)

func GetRequest[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.GetRequest(url))
}

func PostRequest[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.PostRequest(url))
}

func PutRequest[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.PutRequest(url))
}

func DeleteRequest[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.DeleteRequest(url))
}

func Get[RES any](url string, param any) (*RES, error) {
	return NewRequest[RES](http.MethodGet, url).Do(param)
}

func Post[RES any](url string, param any) (*RES, error) {
	return NewRequest[RES](http.MethodPost, url).Do(param)
}

func Put[RES any](url string, param any) (*RES, error) {
	return NewRequest[RES](http.MethodPut, url).Do(param)
}

func Delete[RES any](url string, param any) (*RES, error) {
	return NewRequest[RES](http.MethodDelete, url).Do(param)
}

func GetSubData[RES ResponseInterface[T], T any](url string, param any) (T, error) {
	return NewSubDataRequest[RES, T](client.GetRequest(url)).SubData(param)
}

func GetWithOption[RES ResponseInterface[T], T any](url string, param any, options ...client.Option) (T, error) {
	var response RES
	req := new(client.Client)
	for _, opt := range options {
		opt(req)
	}
	err := req.Get(url, param, &response)
	if err != nil {
		return response.SubData(), err
	}
	return response.SubData(), nil

}
