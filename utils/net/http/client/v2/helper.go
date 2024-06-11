package client

import "github.com/hopeio/cherry/utils/net/http/client"

func NewFromRequest[RES any](req *client.Request) *Request[RES] {
	return (*Request[RES])(req)
}

func NewGetRequest[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.NewGetRequest(url))
}

func DoGet[RES any](url string) (*RES, error) {
	return NewGetRequest[RES](url).Do(nil)
}

func GetSubData[RES ResponseInterface[T], T any](url string) (T, error) {
	return NewSubDataRequestParams[RES, T](client.NewGetRequest(url)).Get(url)
}

func OptionGet[RES ResponseInterface[T], T any](options ...client.Option) func(url string) (T, error) {
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
