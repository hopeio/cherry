package fasthttpctx

import (
	"context"
	"github.com/hopeio/cherry/context/reqctx"
	httpi "github.com/hopeio/cherry/utils/net/http"
	fasthttpi "github.com/hopeio/cherry/utils/net/http/fasthttp"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"
)

type Context = reqctx.Context[*fasthttp.RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[*fasthttp.RequestCtx](ctx)
}

func FromRequest(req *fasthttp.RequestCtx) *Context {
	r := &req.Request

	var ctx context.Context
	if r != nil {
		ctx = req
	}

	ctxi := reqctx.NewRequestContext[*fasthttp.RequestCtx](ctx, req)
	setWithReq(ctxi, r)
	return ctxi
}

func setWithReq(c *Context, r *fasthttp.Request) {
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = Device(&r.Header)
	c.Internal = stringsi.BytesToString(r.Header.Peek(httpi.HeaderGrpcInternal))
}

func Device(r *fasthttp.RequestHeader) *reqctx.DeviceInfo {
	return reqctx.Device(stringsi.BytesToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.BytesToString(r.Peek(httpi.HeaderArea)),
		stringsi.BytesToString(r.Peek(httpi.HeaderLocation)),
		stringsi.BytesToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.BytesToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}

type FastHttpContext Context

func (c *FastHttpContext) SetHeader(md metadata.MD) error {
	header := &c.RequestCtx.Response.Header
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

func (c *FastHttpContext) SendHeader(md metadata.MD) error {
	header := &c.RequestCtx.Response.Header
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

func (c *FastHttpContext) WriteHeader(k, v string) error {
	c.RequestCtx.Response.Header.Set(k, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FastHttpContext) SetCookie(v string) error {
	c.RequestCtx.Response.Header.Set(httpi.HeaderSetCookie, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FastHttpContext) SetTrailer(md metadata.MD) error {
	for k, v := range md {
		if len(v) > 0 {
			c.RequestCtx.Response.Header.Set(k, v[0])
		}
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetTrailer(md)
		if err != nil {
			return err
		}
	}
	return nil
}
