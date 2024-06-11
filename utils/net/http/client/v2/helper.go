package client

import "github.com/hopeio/cherry/utils/net/http/client"

func NewFromRequest[RES any](req *client.Request) *Request[RES] {
	return (*Request[RES])(req)
}

func NewGet[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.NewGet(url))
}

func Get[RES any](url string) (*RES, error) {
	return NewGet[RES](url).Do(nil)
}

func GetSubData[RES ResponseInterface[T], T any](url string) (T, error) {
	return NewSubDataRequestParams[RES, T](client.NewGet(url)).Get(url)
}

func GetWithOption[RES ResponseInterface[T], T any](options ...client.Option) func(url string) (T, error) {
	return func(url string) (T, error) {
		var response RES
		req := new(client.Request)
		for _, opt := range options {
			opt(req)
		}
		err := req.Get(url, &response)
		if err != nil {
			return response.GetData(), err
		}
		return response.GetData(), nil
	}
}
