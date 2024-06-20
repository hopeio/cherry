package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	httpi "github.com/hopeio/cherry/utils/net/http"
	url2 "github.com/hopeio/cherry/utils/net/url"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"io"
	"net/http"
	"strings"
	"time"
)

var DefaultClient = New()

type Request struct {
	ctx         context.Context
	Method, Url string
	contentType ContentType
	headers     httpi.Header //请求级请求头
	client      *Client
}

func NewRequest(method, url string) *Request {
	return &Request{
		ctx:    context.Background(),
		Method: method,
		Url:    url,
		client: DefaultClient,
	}
}

func (req *Request) WithClient(c *Client) *Request {
	req.client = c
	req.client.req = req
	return req
}

func (req *Request) SetClient(set func(c *Client)) *Request {
	req.client = New()
	req.client.req = req
	set(req.client)
	return req
}

func (req *Request) Client() *Client {
	req.client = New()
	req.client.req = req
	return req.client
}

func (req *Request) AddHeader(k, v string) *Request {
	req.headers.Set(k, v)
	return req
}

func (req *Request) ContentType(contentType ContentType) *Request {
	req.contentType = contentType
	return req
}

func (req *Request) Context(ctx context.Context) *Request {
	req.ctx = ctx
	return req
}

func (req *Request) DoEmpty() error {
	return req.Do(nil, nil)
}

func (req *Request) DoNoParam(response any) error {
	return req.Do(nil, response)
}

func (req *Request) DoNoResponse(param any) error {
	return req.Do(param, nil)
}

func (req *Request) DoRaw(param any) (RawBytes, error) {
	var raw RawBytes
	err := req.Do(param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func (req *Request) DoStream(param any) (io.ReadCloser, error) {
	var resp *http.Response
	err := req.Do(param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (r *Request) addHeader(request *http.Request, c *Client) string {
	var auth string
	for k, v := range c.header {
		if k == httpi.HeaderAuthorization {
			auth = v[0]
		}
		request.Header.Add(k, v[0])
	}
	for i := 0; i+1 < len(r.headers); i += 2 {
		request.Header.Set(r.headers[i], r.headers[i+1])
		if r.headers[i] == httpi.HeaderAuthorization {
			auth = r.headers[i+1]
		}
	}

	request.Header.Set(httpi.HeaderContentType, r.contentType.String())
	return auth
}

// Do create a HTTP request
// param: 请求参数 目前只支持编码为json 或 Url-encoded
func (r *Request) Do(param, response any) error {
	if r.Method == "" {
		return errors.New("没有设置请求方法")
	}

	if r.Url == "" {
		return errors.New("没有设置url")
	}
	c := r.client
	var body io.Reader
	var reqBody, respBody *Body
	var statusCode, reqTimes int
	var err error
	var auth string
	reqTime := time.Now()
	// 日志记录
	defer func(now time.Time) {
		if c.logLevel == LogLevelInfo || (err != nil && c.logLevel == LogLevelError) {
			c.logger(r.Method, r.Url, auth, reqBody, respBody, statusCode, time.Since(now), err)
		}
	}(reqTime)

	if r.Method == http.MethodGet {
		r.Url = url2.AppendQueryParam(r.Url, param)
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
				if r.contentType == ContentTypeForm {
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
	request, err = http.NewRequestWithContext(r.ctx, r.Method, r.Url, body)
	if err != nil {
		return err
	}

	auth = r.addHeader(request, c)

	var resp *http.Response
Retry:
	if reqTimes > 0 {
		if c.retryInterval != 0 {
			time.Sleep(c.retryInterval)
		}
		if c.retryHandler != nil {
			c.retryHandler(c)
		}
		reqTime = time.Now()
		if reqBody != nil && reqBody.Data != nil {
			request.Body = io.NopCloser(bytes.NewReader(reqBody.Data))
		}
	}
	resp, err = c.httpClient.Do(request)
	reqTimes++
	if err != nil {
		if c.retryTimes == 0 || reqTimes == c.retryTimes {
			return err
		} else {
			if c.logLevel > LogLevelSilent {
				c.logger(r.Method, r.Url, auth, reqBody, respBody, statusCode, time.Since(reqTime), errors.New(err.Error()+";will retry"))
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
			err = errors.New("status:" + resp.Status + " " + stringsi.ToUnicode(msg))
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
	if c.responseHandler != nil {
		var retry bool
		retry, respBytes, err = c.responseHandler(resp)
		resp.Body.Close()

		if retry {
			if c.logLevel > LogLevelSilent {
				c.logger(r.Method, r.Url, auth, reqBody, respBody, statusCode, time.Since(reqTime), err)
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
