package client_generic

import (
	"context"
	"github.com/hopeio/cherry/utils/net/http/client"
	"net/http"
	"time"
)

// Request ...
type Request[RES any] client.Request

func New[RES any]() *Request[RES] {
	return (*Request[RES])(client.New())
}
func NewRequest[RES any](url string, method string) *Request[RES] {
	return (*Request[RES])(client.NewRequest(url, method))
}

func NewByRequest[RES any](req *client.Request) *Request[RES] {
	return (*Request[RES])(req)
}

func NewGetRequest[RES any](url string) *Request[RES] {
	return (*Request[RES])(client.NewGetRequest(url))
}

func (req *Request[RES]) Origin() *client.Request {
	return (*client.Request)(req)
}

// Do create a HTTP request
func (r *Request[RES]) Do(req any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Do(req, response)
	return response, err
}

func (r *Request[RES]) Get(url string) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Get(url, response)
	return response, err
}

func (r *Request[RES]) Post(url string, param any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Post(url, param, response)
	return response, err
}

func (req *Request[RES]) Context(ctx context.Context) *Request[RES] {
	(*client.Request)(req).Context(ctx)
	return req
}

func (req *Request[RES]) Url(url string) *Request[RES] {
	(*client.Request)(req).Url(url)
	return req
}

func (req *Request[RES]) Method(method string) *Request[RES] {
	(*client.Request)(req).Method(method)
	return req
}

func (req *Request[RES]) ContentType(contentType client.ContentType) *Request[RES] {
	(*client.Request)(req).ContentType(contentType)
	return req
}

func (req *Request[RES]) Header(header client.Header) *Request[RES] {
	(*client.Request)(req).Header(header)
	return req
}

func (req *Request[RES]) AddHeader(k, v string) *Request[RES] {
	(*client.Request)(req).AddHeader(k, v)
	return req
}

func (req *Request[RES]) CachedHeader(key string) *Request[RES] {
	(*client.Request)(req).CachedHeader(key)
	return req
}

func (req *Request[RES]) Logger(logger client.LogCallback) *Request[RES] {
	(*client.Request)(req).Logger(logger)
	return req
}

func (req *Request[RES]) DisableLog() *Request[RES] {
	(*client.Request)(req).DisableLog()
	return req
}

func (req *Request[RES]) LogLevel(lvl client.LogLevel) *Request[RES] {
	(*client.Request)(req).LogLevel(lvl)
	return req
}

func (req *Request[RES]) Tag(tag string) *Request[RES] {
	(*client.Request)(req).Tag(tag)
	return req
}

// handler 返回值:是否重试,返回数据,错误
func (req *Request[RES]) ResponseHandler(handler func(response *http.Response) (retry bool, data []byte, err error)) *Request[RES] {
	(*client.Request)(req).ResponseHandler(handler)
	return req
}

func (req *Request[RES]) Timeout(timeout time.Duration) *Request[RES] {
	(*client.Request)(req).Timeout(timeout)
	return req
}

func (req *Request[RES]) Client(c *http.Client) *Request[RES] {
	(*client.Request)(req).Client(c)
	return req
}

func (req *Request[RES]) RetryTimes(retryTimes int) *Request[RES] {
	(*client.Request)(req).RetryTimes(retryTimes)
	return req
}

func (req *Request[RES]) RetryTimesWithInterval(retryTimes int, retryInterval time.Duration) *Request[RES] {
	(*client.Request)(req).RetryTimesWithInterval(retryTimes, retryInterval)
	return req
}

func (req *Request[RES]) RetryHandler(handle func(*client.Request)) *Request[RES] {
	(*client.Request)(req).RetryHandler(handle)
	return req
}

func (req *Request[RES]) SetRequest(handle func(*client.Request)) *Request[RES] {
	handle((*client.Request)(req))
	return req
}

func (req *Request[RES]) Proxy(url string) *Request[RES] {
	(*client.Request)(req).Proxy(url)
	return req
}

type Option[RES any] func(req *Request[RES]) *Request[RES]
