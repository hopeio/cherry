package client

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
func NewRequest[RES any](method, url string) *Request[RES] {
	return (*Request[RES])(client.NewRequest(method, url))
}

func (r *Request[RES]) Origin() *client.Request {
	return (*client.Request)(r)
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

func (r *Request[RES]) Context(ctx context.Context) *Request[RES] {
	(*client.Request)(r).Context(ctx)
	return r
}

func (r *Request[RES]) Url(url string) *Request[RES] {
	(*client.Request)(r).Url(url)
	return r
}

func (r *Request[RES]) Method(method string) *Request[RES] {
	(*client.Request)(r).Method(method)
	return r
}

func (r *Request[RES]) ContentType(contentType client.ContentType) *Request[RES] {
	(*client.Request)(r).ContentType(contentType)
	return r
}

func (r *Request[RES]) Header(header client.Header) *Request[RES] {
	(*client.Request)(r).Header(header)
	return r
}

func (r *Request[RES]) AddHeader(k, v string) *Request[RES] {
	(*client.Request)(r).AddHeader(k, v)
	return r
}

func (r *Request[RES]) Logger(logger client.AccessLog) *Request[RES] {
	(*client.Request)(r).Logger(logger)
	return r
}

func (r *Request[RES]) DisableLog() *Request[RES] {
	(*client.Request)(r).DisableLog()
	return r
}

func (r *Request[RES]) LogLevel(lvl client.LogLevel) *Request[RES] {
	(*client.Request)(r).LogLevel(lvl)
	return r
}

func (r *Request[RES]) Tag(tag string) *Request[RES] {
	(*client.Request)(r).Tag(tag)
	return r
}

// handler 返回值:是否重试,返回数据,错误
func (r *Request[RES]) ResponseHandler(handler func(response *http.Response) (retry bool, data []byte, err error)) *Request[RES] {
	(*client.Request)(r).ResponseHandler(handler)
	return r
}

func (r *Request[RES]) Timeout(timeout time.Duration) *Request[RES] {
	(*client.Request)(r).Timeout(timeout)
	return r
}

func (r *Request[RES]) Client(c *http.Client) *Request[RES] {
	(*client.Request)(r).Client(c)
	return r
}

func (r *Request[RES]) RetryTimes(retryTimes int) *Request[RES] {
	(*client.Request)(r).RetryTimes(retryTimes)
	return r
}

func (r *Request[RES]) RetryTimesWithInterval(retryTimes int, retryInterval time.Duration) *Request[RES] {
	(*client.Request)(r).RetryTimesWithInterval(retryTimes, retryInterval)
	return r
}

func (r *Request[RES]) RetryHandler(handle func(*client.Request)) *Request[RES] {
	(*client.Request)(r).RetryHandler(handle)
	return r
}

func (r *Request[RES]) BasicAuth(authUser, authPass string) {
	(*client.Request)(r).BasicAuth(authUser, authPass)
}

func (r *Request[RES]) SetRequest(handle func(*client.Request)) *Request[RES] {
	handle((*client.Request)(r))
	return r
}

func (r *Request[RES]) Proxy(url string) *Request[RES] {
	(*client.Request)(r).Proxy(url)
	return r
}

type Option[RES any] func(req *Request[RES])
