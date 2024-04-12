package client_generic

import "github.com/hopeio/cherry/utils/net/http/client"

func SimpleGet[RES any](url string) (*RES, error) {
	return NewGetRequest[RES](url).Do(nil)
}

func GetSubData[RES ResponseInterface[T], T any](url string) (T, error) {
	return NewSubDataRequestParams[RES, T](client.NewGetRequest(url)).Get(url)
}

// Deprecated
func CustomizeGet[RES ResponseInterface[T], T any](option client.Option) func(url string) (T, error) {
	return func(url string) (T, error) {
		var response RES
		err := option(new(client.Request)).Get(url, &response)
		if err != nil {
			return response.GetData(), err
		}
		return response.GetData(), nil
	}
}
