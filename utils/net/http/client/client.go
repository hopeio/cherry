package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	httpi "github.com/hopeio/cherry/utils/net/http"
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
	defaultClient = newHttpClient()
	globalLevel   = LogLevelError
	headerMap     = sync.Map{}
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
	defaultClient.Timeout = timeout
}

func DisableLog() {
	globalLevel = LogLevelSilent
}

func SetDefaultLogger(logger LogCallback) {
	defaultLog = logger
}

func SetProxy(url string) {
	purl, _ := stdurl.Parse(url)
	setProxy(defaultClient, http.ProxyURL(purl))
}

func ResetProxy() {
	defaultClient.Transport.(*http.Transport).Proxy = http.ProxyFromEnvironment
}

func SetDefaultClient(client *http.Client) {
	defaultClient = client
}

type ContentType uint8

const (
	ContentTypeJson ContentType = iota
	ContentTypeForm
	ContentTypeFormData
	ContentTypeProtobuf
	ContentTypeText
	ContentTypeImage
	ContentTypeBinary
)

// Request ...
type Request struct {
	ctx context.Context
	// client settings
	client *http.Client
	// 适用于单次配置不同的请求,如果设置是固定,建议设置0值，直接设置client
	timeout time.Duration
	// 内部使用标志性字段,用于判断是否重复设置代理,为空时代表从环境变量获取
	clientProxy string
	proxyUrl    string
	tag         string // 默认json

	// request
	url, method        string
	contentType        ContentType
	AuthUser, AuthPass string
	header             Header
	cachedHeaderKey    string

	// response
	responseHandler func(response *http.Response) (retry bool, data []byte, err error)

	// config
	logger   LogCallback
	logLevel LogLevel

	// retry
	retryTimes    int
	retryInterval time.Duration
	retryHandler  func(*Request)
}

func New() *Request {
	return newRequest("", "")
}

func NewRequest(url, method string) *Request {
	return newRequest(url, strings.ToUpper(method))
}

func newRequest(url, method string) *Request {
	return &Request{ctx: context.Background(), client: defaultClient, url: url, method: method, header: make([]string, 0, 2), logger: defaultLog, logLevel: globalLevel, retryInterval: 200 * time.Millisecond}
}

func DefaultHeaderRequest() *Request {
	req := newRequest("", "")
	req.Header(DefaultHeader())
	return req
}

func NewGetRequest(url string) *Request {
	return newRequest(url, http.MethodGet)
}

func NewPostRequest(url string) *Request {
	return newRequest(url, http.MethodPost)
}

func NewPutRequest(url string) *Request {
	return newRequest(url, http.MethodPut)
}

func NewDeleteRequest(url string) *Request {
	return newRequest(url, http.MethodDelete)
}

func (req *Request) Context(ctx context.Context) *Request {
	req.ctx = ctx
	return req
}

func (req *Request) Url(url string) *Request {
	req.url = url
	return req
}

func (req *Request) Method(method string) *Request {
	req.method = strings.ToUpper(method)
	return req
}

func (req *Request) ContentType(contentType ContentType) *Request {
	req.contentType = contentType
	return req
}

func (req *Request) Header(header Header) *Request {
	req.header = header
	return req
}

func (req *Request) AddHeader(k, v string) *Request {
	req.header = append(req.header, k, v)
	return req
}

func (req *Request) CachedHeader(key string) *Request {
	req.cachedHeaderKey = key
	return req
}

func (req *Request) Logger(logger LogCallback) *Request {
	if logger == nil {
		return req
	}
	req.logger = logger
	return req
}

func (req *Request) DisableLog() *Request {
	req.logLevel = LogLevelSilent
	return req
}

func (req *Request) LogLevel(lvl LogLevel) *Request {
	req.logLevel = lvl
	return req
}

func (req *Request) Tag(tag string) *Request {
	req.tag = tag
	return req
}

// handler 返回值:是否重试,返回数据,错误
func (req *Request) ResponseHandler(handler func(response *http.Response) (retry bool, data []byte, err error)) *Request {
	req.responseHandler = handler
	return req
}

// 设置过期时间,仅对单次请求有效
func (req *Request) Timeout(timeout time.Duration) *Request {
	req.timeout = timeout
	return req
}

func (req *Request) Client(client *http.Client) *Request {
	req.client = client
	return req
}

func (req *Request) RetryTimes(retryTimes int) *Request {
	req.retryTimes = retryTimes
	return req
}

func (req *Request) RetryTimesWithInterval(retryTimes int, retryInterval time.Duration) *Request {
	req.retryTimes = retryTimes
	req.retryInterval = retryInterval
	return req
}

func (req *Request) RetryHandler(handle func(*Request)) *Request {
	req.retryHandler = handle
	return req
}

func (req *Request) Proxy(url string) *Request {
	req.proxyUrl = url
	return req
}

func (req *Request) BasicAuth(authUser, authPass string) {

}

type ResponseBodyCheck interface {
	CheckError() error
}

type RawBytes = []byte

func (req *Request) DoWithNoParam(response interface{}) error {
	return req.Do(nil, response)
}

func (req *Request) DoWithNoResponse(param interface{}) error {
	return req.Do(param, nil)
}

func (req *Request) DoEmpty() error {
	return req.Do(nil, nil)
}

func (req *Request) setHeader(request *http.Request) {
	for i := 0; i+1 < len(req.header); i += 2 {
		request.Header.Set(req.header[i], req.header[i+1])
	}
	if req.AuthUser != "" && req.AuthPass != "" {
		request.SetBasicAuth(req.AuthUser, req.AuthPass)
	}
	if req.contentType == ContentTypeJson {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentJsonHeaderValue)
	} else if req.contentType == ContentTypeFormData {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentFormHeaderValue)
	} else {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentFormMultipartHeaderValue)
	}
}

// Do create a HTTP request
func (req *Request) Do(param, response interface{}) error {
	if req.method == "" {
		return errors.New("没有设置请求方法")
	}
	method := req.method
	if req.url == "" {
		return errors.New("没有设置url")
	}
	url := req.url
	if req.client == nil {
		req.client = defaultClient
	}
	if req.timeout != 0 && req.timeout != req.client.Timeout {
		defer setTimeout(req.client, req.client.Timeout)
		setTimeout(req.client, req.timeout)
	}
	if req.proxyUrl != "" && req.proxyUrl != req.clientProxy {
		purl, _ := stdurl.Parse(url)
		setProxy(req.client, http.ProxyURL(purl))
		req.clientProxy = url
	} else if req.clientProxy != "" {
		setProxy(req.client, http.ProxyFromEnvironment)
		req.clientProxy = ""
	}
	var body io.Reader
	var reqBody, respBody *Body
	var statusCode, reqTimes int
	var err error
	reqTime := time.Now()
	// 日志记录
	defer func(now time.Time) {
		if req.logLevel == LogLevelInfo || (err != nil && req.logLevel == LogLevelError) {
			req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Since(now), err)
		}
	}(reqTime)

	if method == http.MethodGet {
		url = UrlAppendParam(req.url, param)
	} else {
		reqBody = &Body{}
		if param != nil {
			switch paramt := param.(type) {
			case string:
				body = strings.NewReader(paramt)
				reqBody.Data = stringsi.ToBytes(paramt)
			case []byte:
				body = bytes.NewReader(paramt)
				reqBody.Data = paramt
			case io.Reader:
				var reqBytes []byte
				reqBytes, err = io.ReadAll(paramt)
				body = bytes.NewReader(reqBytes)
				reqBody.Data = reqBytes
			default:
				if req.contentType == ContentTypeJson {
					var reqBytes []byte
					reqBytes, err = json.Marshal(param)
					if err != nil {
						return err
					}
					body = bytes.NewReader(reqBytes)
					reqBody.Data = reqBytes
					reqBody.ContentType = ContentTypeJson
				} else {
					params := UrlParam(param)
					reqBody.Data = stringsi.ToBytes(params)
					body = strings.NewReader(params)
				}
			}
		}
	}
	var request *http.Request
	request, err = http.NewRequestWithContext(req.ctx, method, url, body)
	if err != nil {
		return err
	}

	// 缓存header
	if req.cachedHeaderKey != "" {
		if header, ok := headerMap.Load(req.cachedHeaderKey); ok {
			request.Header = header.(http.Header)
		} else {
			req.setHeader(request)
			headerMap.Store(req.cachedHeaderKey, request.Header)
		}
	} else {
		req.setHeader(request)
	}

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
				req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Since(reqTime), errors.New(err.Error()+";will retry"))
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
	if resp.Header.Get(httpi.HeaderContentEncoding) == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return err
		}
	} else {
		reader = resp.Body
	}

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
		if err != nil {
			return err
		}

		if retry {
			if req.logLevel > LogLevelSilent {
				req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Since(reqTime), errors.New("will retry"))
			}
			goto Retry
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
		if strings.HasPrefix(contentType, httpi.ContentJsonHeaderValue) {
			respBody.ContentType = ContentTypeJson
		} else if strings.HasPrefix(contentType, httpi.ContentFormHeaderValue) {
			respBody.ContentType = ContentTypeForm
		} else if strings.HasPrefix(contentType, "text") {
			respBody.ContentType = ContentTypeText
		} else if strings.HasPrefix(contentType, "image") {
			respBody.ContentType = ContentTypeImage
		} else {
			respBody.ContentType = ContentTypeJson
		}

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
				return errors.New("json.Unmarshal error:" + err.Error())
			}
		}

		if v, ok := response.(ResponseBodyCheck); ok {
			err = v.CheckError()
		}
	}

	return err
}

func (req *Request) DoRaw(param interface{}) (RawBytes, error) {
	var raw RawBytes
	err := req.Do(param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func (req *Request) DoStream(param interface{}) (io.ReadCloser, error) {
	var resp *http.Response
	err := req.Do(param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (req *Request) Get(url string, response interface{}) error {
	req.url = url
	req.method = http.MethodGet
	return req.Do(nil, response)
}

func (req *Request) Post(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPost
	return (req).Do(param, response)
}

func (req *Request) Put(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPut
	return req.Do(param, response)
}

func (req *Request) Delete(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodDelete
	return req.Do(param, response)
}
