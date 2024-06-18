package client

import (
	"context"
	"github.com/hopeio/cherry/utils/net/http/client"
	"net/http"
	"time"
)

// Client ...
type Client[RES any] client.Client

func New[RES any]() *Client[RES] {
	return (*Client[RES])(client.New())
}

func (r *Client[RES]) Origin() *client.Client {
	return (*client.Client)(r)
}

// Do create a HTTP request
func (r *Client[RES]) Do(method, url string, req any) (*RES, error) {
	response := new(RES)
	err := (*client.Client)(r).Do(method, url, req, response)
	return response, err
}

func (r *Client[RES]) Get(url string, param any) (*RES, error) {
	response := new(RES)
	err := (*client.Client)(r).Get(url, param, response)
	return response, err
}

func (r *Client[RES]) Post(url string, param any) (*RES, error) {
	response := new(RES)
	err := (*client.Client)(r).Post(url, param, response)
	return response, err
}

func (r *Client[RES]) Context(ctx context.Context) *Client[RES] {
	(*client.Client)(r).Context(ctx)
	return r
}

func (r *Client[RES]) ContentType(contentType client.ContentType) *Client[RES] {
	(*client.Client)(r).ContentType(contentType)
	return r
}

func (r *Client[RES]) Header(header http.Header) *Client[RES] {
	(*client.Client)(r).Header(header)
	return r
}

func (r *Client[RES]) AddHeader(k, v string) *Client[RES] {
	(*client.Client)(r).AddHeader(k, v)
	return r
}

func (r *Client[RES]) Logger(logger client.AccessLog) *Client[RES] {
	(*client.Client)(r).Logger(logger)
	return r
}

func (r *Client[RES]) DisableLog() *Client[RES] {
	(*client.Client)(r).DisableLog()
	return r
}

func (r *Client[RES]) LogLevel(lvl client.LogLevel) *Client[RES] {
	(*client.Client)(r).LogLevel(lvl)
	return r
}

func (r *Client[RES]) ParseTag(tag string) *Client[RES] {
	(*client.Client)(r).ParseTag(tag)
	return r
}

// handler 返回值:是否重试,返回数据,错误
func (r *Client[RES]) ResponseHandler(handler func(response *http.Response) (retry bool, data []byte, err error)) *Client[RES] {
	(*client.Client)(r).ResponseHandler(handler)
	return r
}

func (r *Client[RES]) Timeout(timeout time.Duration) *Client[RES] {
	(*client.Client)(r).Timeout(timeout)
	return r
}

func (r *Client[RES]) HttpClient(c *http.Client) *Client[RES] {
	(*client.Client)(r).HttpClient(c)
	return r
}

func (r *Client[RES]) RetryTimes(retryTimes int) *Client[RES] {
	(*client.Client)(r).RetryTimes(retryTimes)
	return r
}

func (r *Client[RES]) RetryTimesWithInterval(retryTimes int, retryInterval time.Duration) *Client[RES] {
	(*client.Client)(r).RetryTimesWithInterval(retryTimes, retryInterval)
	return r
}

func (r *Client[RES]) RetryHandler(handle func(*client.Client)) *Client[RES] {
	(*client.Client)(r).RetryHandler(handle)
	return r
}

func (r *Client[RES]) BasicAuth(authUser, authPass string) {
	(*client.Client)(r).BasicAuth(authUser, authPass)
}

func (r *Client[RES]) SetRequest(handle func(*client.Client)) *Client[RES] {
	handle((*client.Client)(r))
	return r
}

func (r *Client[RES]) Proxy(url string) *Client[RES] {
	(*client.Client)(r).Proxy(url)
	return r
}

type Option[RES any] func(req *Client[RES])
