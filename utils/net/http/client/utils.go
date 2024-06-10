package client

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func SetProxyEnv(url string) {
	os.Setenv("HTTP_PROXY", url)
	os.Setenv("HTTPS_PROXY", url)
}

func ResolveURL(u *url.URL, p string) string {
	if strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://") {
		return p
	}
	var baseURL string
	if strings.Index(p, "/") == 0 {
		baseURL = u.Scheme + "://" + u.Host
	} else {
		tU := u.String()
		baseURL = tU[0:strings.LastIndex(tU, "/")]
	}
	return baseURL + path.Join("/", p)
}

func setTimeout(client *http.Client, timeout time.Duration) {
	if client == nil {
		client = DefaultClient
	}
	if timeout < time.Second {
		timeout = timeout * time.Second
	}
	client.Timeout = timeout
}

func setProxy(client *http.Client, proxy func(*http.Request) (*url.URL, error)) {
	client.Transport.(*http.Transport).Proxy = proxy
}
