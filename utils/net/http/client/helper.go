package client

import (
	"io"
	"net/http"
)

func SimpleGet(url string, response any) error {
	return NewGetRequest(url).DisableLog().DoWithNoParam(response)
}

func SimpleGetStream(url string) (io.ReadCloser, error) {
	var resp *http.Response
	err := SimpleGet(url, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func SimplePost(url string, param, response interface{}) error {
	return NewPostRequest(url).DisableLog().Do(param, response)
}

func SimplePut(url string, param, response interface{}) error {
	return NewPutRequest(url).DisableLog().Do(param, response)
}

func SimpleDelete(url string, param, response interface{}) error {
	return NewDeleteRequest(url).DisableLog().Do(param, response)
}
