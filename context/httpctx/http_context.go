package httpctx

import (
	"context"
	"github.com/hopeio/cherry/context/reqctx"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type RequestCtx struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type Context = reqctx.Context[RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[RequestCtx](ctx)
}

func FromRequest(req RequestCtx) *Context {
	r := req.Request

	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	}

	ctxi := reqctx.NewRequestContext[RequestCtx](ctx, req)
	setWithHttpReq(ctxi, r)
	return ctxi
}

func setWithHttpReq(c *reqctx.Context[RequestCtx], r *http.Request) {
	if r == nil {
		return
	}
	c.DeviceInfo = DeviceFromHeader(r.Header)
	c.Internal = r.Header.Get(httpi.HeaderGrpcInternal)
	c.Token = httpi.GetToken(r)
}

func DeviceFromHeader(r http.Header) *reqctx.DeviceInfo {
	return reqctx.Device(r.Get(httpi.HeaderDeviceInfo),
		r.Get(httpi.HeaderArea), r.Get(httpi.HeaderLocation),
		r.Get(httpi.HeaderUserAgent), r.Get(httpi.HeaderXForwardedFor))
}

type ReqValue[REQ any, V any] struct {
	reqctx.ReqValue[REQ]
	Value V
}

type HttpContext Context

func (c *HttpContext) SetHeader(md metadata.MD) error {
	header := c.RequestCtx.Response.Header()
	for k, v := range md {
		if len(v) > 0 {
			header.Set(k, v[0])
		}
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *HttpContext) SendHeader(md metadata.MD) error {
	header := c.RequestCtx.Response.Header()
	for k, v := range md {
		if len(v) > 0 {
			header.Set(k, v[0])
		}
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *HttpContext) WriteHeader(k, v string) error {
	c.RequestCtx.Response.Header().Set(k, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *HttpContext) SetCookie(v string) error {
	c.RequestCtx.Response.Header().Set(httpi.HeaderSetCookie, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *HttpContext) SetTrailer(md metadata.MD) error {
	for k, v := range md {
		c.RequestCtx.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetTrailer(md)
		if err != nil {
			return err
		}
	}
	return nil
}
