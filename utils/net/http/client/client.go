package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	httpi "github.com/hopeio/cherry/utils/net/http"
	url2 "github.com/hopeio/cherry/utils/net/url"
	"io"
	"net"
	"net/http"
	stdurl "net/url"
	"strings"
	"sync"
	"time"

	stringsi "github.com/hopeio/cherry/utils/strings"
)

// 不是并发安全的

var (
	DefaultHttpClient = newHttpClient()
	DefaultLogLevel   = LogLevelError
	headerMap         = sync.Map{}
)

var timeout = time.Minute

func newHttpClient() *http.Client {
	return &http.Client{
		//Timeout: timeout * 2,
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment, // 代理使用
			ForceAttemptHTTP2: true,
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			//DisableKeepAlives: true,
			TLSHandshakeTimeout: timeout,
		},
	}
}

func SetTimeout(timeout time.Duration) {
	DefaultHttpClient.Timeout = timeout
}

func DisableLog() {
	DefaultLogLevel = LogLevelSilent
}

func SetAccessLog(log AccessLog) {
	defaultLog = log
}

func SetProxy(url string) {
	purl, _ := stdurl.Parse(url)
	setProxy(DefaultHttpClient, http.ProxyURL(purl))
}

func ResetProxy() {
	DefaultHttpClient.Transport.(*http.Transport).Proxy = http.ProxyFromEnvironment
}

func SetHttpClient(client *http.Client) {
	DefaultHttpClient = client
}

// Client ...
type Client struct {
	ctx context.Context
	// client settings
	client        *http.Client
	defaultClient bool
	// 适用于单次配置不同的请求,如果设置是固定,建议设置0值，直接设置client
	timeout time.Duration

	proxyUrl string
	parseTag string // 默认json

	// request
	contentType        ContentType
	authUser, authPass string
	header             http.Header

	// response
	responseHandler func(response *http.Response) (retry bool, data []byte, err error)

	// logger
	logger   AccessLog
	logLevel LogLevel

	// retry
	retryTimes    int
	retryInterval time.Duration
	retryHandler  func(*Client)
}

func New() *Client {
	return newClient()
}

func newClient() *Client {
	return &Client{ctx: context.Background(), client: DefaultHttpClient, logger: defaultLog, logLevel: DefaultLogLevel, retryInterval: 200 * time.Millisecond}
}

func (req *Client) Context(ctx context.Context) *Client {
	req.ctx = ctx
	return req
}

func (req *Client) ContentType(contentType ContentType) *Client {
	req.contentType = contentType
	return req
}

func (req *Client) Header(header http.Header) *Client {
	req.header = header
	return req
}

func (req *Client) AddHeader(k, v string) *Client {
	if req.header == nil {
		req.header = make(http.Header)
	}
	req.header.Set(k, v)
	return req
}

func (req *Client) Logger(logger AccessLog) *Client {
	if logger == nil {
		return req
	}
	req.logger = logger
	return req
}

func (req *Client) DisableLog() *Client {
	req.logLevel = LogLevelSilent
	return req
}

func (req *Client) LogLevel(lvl LogLevel) *Client {
	req.logLevel = lvl
	return req
}

func (req *Client) ParseTag(tag string) *Client {
	req.parseTag = tag
	return req
}

// handler 返回值:是否重试,返回数据,错误
func (req *Client) ResponseHandler(handler func(response *http.Response) (retry bool, data []byte, err error)) *Client {
	req.responseHandler = handler
	return req
}

// 设置过期时间,仅对单次请求有效
func (req *Client) Timeout(timeout time.Duration) *Client {
	req.timeout = timeout
	return req
}

func (req *Client) HttpClient(client *http.Client) *Client {
	req.client = client
	return req
}

func (req *Client) RetryTimes(retryTimes int) *Client {
	req.retryTimes = retryTimes
	return req
}

func (req *Client) RetryTimesWithInterval(retryTimes int, retryInterval time.Duration) *Client {
	req.retryTimes = retryTimes
	req.retryInterval = retryInterval
	return req
}

func (req *Client) RetryHandler(handle func(*Client)) *Client {
	req.retryHandler = handle
	return req
}

func (req *Client) Proxy(url string) *Client {
	req.proxyUrl = url
	return req
}

func (req *Client) BasicAuth(authUser, authPass string) *Client {
	req.authUser, req.authPass = authUser, authPass
	return req
}

func (req *Client) Request(method, url string) *Request {
	return &Request{
		method: method, url: url, Client: req,
	}
}

func (req *Client) Clone() *Client {
	return &(*req)
}

type ResponseBodyCheck interface {
	CheckError() error
}

type RawBytes = []byte

func (req *Client) DoNoParam(method, url string, response interface{}) error {
	return req.Do(method, url, nil, response)
}

func (req *Client) DoNoResponse(method, url string, param interface{}) error {
	return req.Do(method, url, param, nil)
}

func (req *Client) DoEmpty(method, url string) error {
	return req.Do(method, url, nil, nil)
}

func (req *Client) addHeader(request *http.Request) {
	if req.authUser != "" && req.authPass != "" {
		request.SetBasicAuth(req.authUser, req.authPass)
	}
	request.Header.Set(httpi.HeaderContentType, req.contentType.String())
}

// Do create a HTTP request
// param: 请求参数 目前只支持编码为json 或 url-encoded
func (req *Client) Do(method, url string, param, response interface{}) error {
	if method == "" {
		return errors.New("没有设置请求方法")
	}

	if url == "" {
		return errors.New("没有设置url")
	}

	if req.client == nil {
		req.client = DefaultHttpClient
		req.defaultClient = true
	}
	if req.timeout != 0 && req.timeout != req.client.Timeout {
		if req.defaultClient {
			req.client = newHttpClient()
			req.defaultClient = false
		}
		setTimeout(req.client, req.timeout)
	}
	if req.proxyUrl != "" {
		purl, _ := req.client.Transport.(*http.Transport).Proxy(nil)
		if req.proxyUrl != purl.String() {
			if req.defaultClient {
				req.client = newHttpClient()
				req.defaultClient = false
			}
			purl, _ = stdurl.Parse(url)
			setProxy(req.client, http.ProxyURL(purl))
		}
	}
	var body io.Reader
	var reqBody, respBody *Body
	var statusCode, reqTimes int
	var err error
	reqTime := time.Now()
	// 日志记录
	defer func(now time.Time) {
		if req.logLevel == LogLevelInfo || (err != nil && req.logLevel == LogLevelError) {
			req.logger(method, url, req.authUser, reqBody, respBody, statusCode, time.Since(now), err)
		}
	}(reqTime)

	if method == http.MethodGet {
		url = url2.AppendQueryParam(url, param)
	} else {
		reqBody = &Body{}
		if param != nil {
			switch paramType := param.(type) {
			case string:
				body = strings.NewReader(paramType)
				reqBody.Data = stringsi.ToBytes(paramType)
			case []byte:
				body = bytes.NewReader(paramType)
				reqBody.Data = paramType
			case io.Reader:
				var reqBytes []byte
				reqBytes, err = io.ReadAll(paramType)
				body = bytes.NewReader(reqBytes)
				reqBody.Data = reqBytes
			default:
				if req.contentType == ContentTypeForm {
					params := url2.QueryParam(param)
					reqBody.Data = stringsi.ToBytes(params)
					body = strings.NewReader(params)
				} else {
					var reqBytes []byte
					reqBytes, err = json.Marshal(param)
					if err != nil {
						return err
					}
					body = bytes.NewReader(reqBytes)
					reqBody.Data = reqBytes
					reqBody.ContentType = ContentTypeJson
				}
			}
		}
	}
	var request *http.Request
	request, err = http.NewRequestWithContext(req.ctx, method, url, body)
	if err != nil {
		return err
	}

	req.addHeader(request)

	var resp *http.Response
Retry:
	if reqTimes > 0 {
		if req.retryInterval != 0 {
			time.Sleep(req.retryInterval)
		}
		if req.retryHandler != nil {
			req.retryHandler(req)
		}
		reqTime = time.Now()
		if reqBody != nil && reqBody.Data != nil {
			request.Body = io.NopCloser(bytes.NewReader(reqBody.Data))
		}
	}
	resp, err = req.client.Do(request)
	reqTimes++
	if err != nil {
		if req.retryTimes == 0 || reqTimes == req.retryTimes {
			return err
		} else {
			if req.logLevel > LogLevelSilent {
				req.logger(method, url, req.authUser, reqBody, respBody, statusCode, time.Since(reqTime), errors.New(err.Error()+";will retry"))
			}
			goto Retry
		}
	}

	respBody = &Body{}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		respBody.ContentType = ContentTypeText
		if resp.StatusCode == http.StatusNotFound {
			err = errors.New("not found")
		} else {
			var msg []byte
			msg, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}
			err = errors.New("status:" + resp.Status + " " + stringsi.ConvertUnicode(msg))
		}
		return err
	}

	if httpresp, ok := response.(*http.Response); ok {
		*httpresp = *resp
		return err
	}

	if httpresp, ok := response.(**http.Response); ok {
		*httpresp = resp
		return err
	}

	var reader io.Reader
	// net/http会自动处理gzip
	/*	if resp.Header.Get(httpi.HeaderContentEncoding) == "gzip" {
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				resp.Body.Close()
				return err
			}
		} else {
			reader = resp.Body
		}*/

	reader = resp.Body

	if httpresp, ok := response.(*io.Reader); ok {
		*httpresp = reader
		return err
	}
	statusCode = resp.StatusCode

	var respBytes []byte
	if req.responseHandler != nil {
		var retry bool
		retry, respBytes, err = req.responseHandler(resp)
		resp.Body.Close()

		if retry {
			if req.logLevel > LogLevelSilent {
				req.logger(method, url, req.authUser, reqBody, respBody, statusCode, time.Since(reqTime), err)
			}
			goto Retry
		} else if err != nil {
			return err
		}
	} else {
		respBytes, err = io.ReadAll(reader)
		resp.Body.Close()
		if err != nil {
			return err
		}
	}
	respBody.Data = respBytes
	if len(respBytes) > 0 && response != nil {
		contentType := resp.Header.Get(httpi.HeaderContentType)
		respBody.ContentType.Decode(contentType)

		if raw, ok := response.(*RawBytes); ok {
			*raw = respBytes
			return nil
		}
		if respBody.ContentType == ContentTypeForm {
			// TODO
		} else {
			// 默认json
			err = json.Unmarshal(respBytes, response)
			if err != nil {
				return fmt.Errorf("json.Unmarshal error: %v", err)
			}
		}

		if v, ok := response.(ResponseBodyCheck); ok {
			err = v.CheckError()
		}
	}

	return err
}

func (req *Client) DoRaw(method, url string, param interface{}) (RawBytes, error) {
	var raw RawBytes
	err := req.Do(method, url, param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func (req *Client) DoStream(method, url string, param interface{}) (io.ReadCloser, error) {
	var resp *http.Response
	err := req.Do(method, url, param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (req *Client) Get(url string, param, response interface{}) error {
	return req.Do(http.MethodGet, url, param, response)
}

func (req *Client) GetNP(url string, response interface{}) error {
	return req.Do(http.MethodGet, url, nil, response)
}

func (req *Client) Post(url string, param, response interface{}) error {
	return req.Do(http.MethodPost, url, param, response)
}

func (req *Client) Put(url string, param, response interface{}) error {
	return req.Do(http.MethodPut, url, param, response)
}

func (req *Client) Delete(url string, param, response interface{}) error {
	return req.Do(http.MethodDelete, url, param, response)
}
