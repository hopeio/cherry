package client

import "github.com/hopeio/cherry/utils/net/http/client"

func SimpleGet[RES any](url string) (*RES, error) {
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
