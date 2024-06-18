package client

import (
	httpi "github.com/hopeio/cherry/utils/net/http"
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
