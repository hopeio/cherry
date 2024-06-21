package client

import (
	httpi "github.com/hopeio/cherry/utils/net/http"
	"io"
	"net/http"
)

func DefaultHeader() http.Header {
	return http.Header{
		httpi.HeaderAcceptLanguage: []string{"zh-CN,zh;q=0.9;charset=utf-8"},
		httpi.HeaderConnection:     []string{"keep-alive"},
		httpi.HeaderUserAgent:      []string{UserAgentChrome117},
		//"Accept", "application/json,text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8", // 将会越来越少用，服务端一般固定格式
	}
}

func DefaultHeaderClient() *Client {
	return New().Header(DefaultHeader())
}

func DefaultHeaderRequest() *Request {
	return &Request{client: New().Header(DefaultHeader())}
}

func GetRequest(url string) *Request {
	return NewRequest(http.MethodGet, url)
}

func PostRequest(url string) *Request {
	return NewRequest(http.MethodPost, url)
}

func PutRequest(url string) *Request {
	return NewRequest(http.MethodPut, url)
}

func DeleteRequest(url string) *Request {
	return NewRequest(http.MethodDelete, url)
}

func Get(url string, param, response any) error {
	return GetRequest(url).Do(param, response)
}

func GetX(url string, response any) error {
	return Get(url, nil, response)
}

func GetStream(url string, param any) (io.ReadCloser, error) {
	var resp *http.Response
	err := Get(url, param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func GetStreamX(url string) (io.ReadCloser, error) {
	return GetStream(url, nil)
}

func Post(url string, param, response interface{}) error {
	return New().DisableLog().Post(url, param, response)
}

func Put(url string, param, response interface{}) error {
	return New().DisableLog().Put(url, param, response)
}

func Delete(url string, param, response interface{}) error {
	return New().DisableLog().Delete(url, param, response)
}
