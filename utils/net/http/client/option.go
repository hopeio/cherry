package client

import (
	httpi "github.com/hopeio/cherry/utils/net/http"
	"net/http"
)

type Option func(req *Request) *Request

type RequestOption func(req *http.Request)

func SetHeader(header Header) RequestOption {
	return func(req *http.Request) {
		for i := 0; i < len(header)-1; i += 2 {
			req.Header.Set(header[i], header[i+1])
		}

	}
}

func SetRefer(refer string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderReferer, refer)
	}
}

func SetAccept(refer string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderAccept, refer)
	}
}

func SetCookie(cookie string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderCookie, cookie)
	}
}

// TODO
// tag :`request:"uri:xxx;query:xxx;header:xxx;body:xxx"`
func setRequest(p any, req *http.Request) {

}
