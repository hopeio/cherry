package client

import "github.com/hopeio/cherry/utils/net/http/client"

func NewFromClient[RES any](req *client.Client) *Client[RES] {
	return (*Client[RES])(req)
}

func NewGet[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.NewGet(url))
}

func Get[RES any](url string, param any) (*RES, error) {
	return (*Request[RES])(client.NewGet(url)).Do(param)
}

func GetSubData[RES ResponseInterface[T], T any](url string) (T, error) {
	return NewSubDataRequestParams[RES, T](client.NewGet(url)).Get(url)
}

func GetWithOption[RES ResponseInterface[T], T any](url string, options ...client.Option) (T, error) {
	var response RES
	req := new(client.Client)
	for _, opt := range options {
		opt(req)
	}
	err := req.Get(url, nil, &response)
	if err != nil {
		return response.GetData(), err
	}
	return response.GetData(), nil

}
