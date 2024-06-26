package client

import (
	httpi "github.com/hopeio/cherry/utils/net/http"
	"net/http"
)

type Option func(req *Client) *Client

type HttpRequestOption func(req *http.Request)

func AddHeader(k, v string) HttpRequestOption {
	return func(req *http.Request) {
		req.Header.Set(k, v)
	}
}

func SetRefer(refer string) HttpRequestOption {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderReferer, refer)
	}
}

func SetAccept(refer string) HttpRequestOption {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderAccept, refer)
	}
}

func SetCookie(cookie string) HttpRequestOption {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderCookie, cookie)
	}
}

// TODO
// tag :`request:"uri:xxx;query:xxx;header:xxx;body:xxx"`
func setRequest(p any, req *http.Request) {

}

type HttpClientOption func(client *http.Client)
type ResponseHandler func(client *http.Response)
