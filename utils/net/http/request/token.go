package request

import (
	httpi "github.com/hopeio/cherry/utils/net/http"
	"net/http"
	"net/url"
)

func GetToken(r *http.Request) string {
	if token := r.Header.Get(httpi.HeaderAuthorization); token != "" {
		return token
	}
	if cookie, _ := r.Cookie(httpi.HeaderCookieToken); cookie != nil {
		value, _ := url.QueryUnescape(cookie.Value)
		return value
	}
	return ""
}
