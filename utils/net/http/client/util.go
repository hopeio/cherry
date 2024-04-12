package client

import (
	"github.com/hopeio/cherry/utils/number"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func UrlParam(param interface{}) string {
	if param == nil {
		return ""
	}
	query := url.Values{}
	parseParam(param, query)
	return query.Encode()
}

func parseParam(param interface{}, query url.Values) {
	v := reflect.ValueOf(param).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		filed := v.Field(i)
		kind := filed.Kind()
		if kind == reflect.Interface || kind == reflect.Ptr {
			parseParam(filed.Interface(), query)
			continue
		}
		if kind == reflect.Struct {
			parseParam(filed.Addr().Interface(), query)
			continue
		}
		value := getFieldValue(filed)
		if value != "" {
			query.Set(t.Field(i).Tag.Get("json"), getFieldValue(v.Field(i)))
		}
	}

}

func getFieldValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.Itoa(int(v.Int()))
	case reflect.Float32, reflect.Float64:
		return number.FormatFloat(v.Float())
	case reflect.String:
		return v.String()
	case reflect.Interface, reflect.Ptr:
		return getFieldValue(v.Elem())
	case reflect.Struct:

	}
	return ""
}

func UrlAppendParam(url string, param interface{}) string {
	if param == nil {
		return url
	}
	sep := "?"
	if strings.Contains(url, sep) {
		sep = "&"
	}
	switch paramt := param.(type) {
	case string:
		url += sep + paramt
	case []byte:
		url += sep + stringsi.BytesToString(paramt)
	default:
		params := UrlParam(param)
		url += sep + params
	}
	return url
}

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
		client = defaultClient
	}
	if timeout < time.Second {
		timeout = timeout * time.Second
	}
	client.Timeout = timeout
}

func setProxy(client *http.Client, proxy func(*http.Request) (*url.URL, error)) {
	client.Transport.(*http.Transport).Proxy = proxy
}
