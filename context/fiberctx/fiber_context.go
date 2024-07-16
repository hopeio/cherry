package fiberctx

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/context/fasthttpctx"
	"github.com/hopeio/cherry/context/reqctx"
	httpi "github.com/hopeio/cherry/utils/net/http"
	fasthttpi "github.com/hopeio/cherry/utils/net/http/fasthttp"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type Context = reqctx.Context[fiber.Ctx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[fiber.Ctx](ctx)
}

func FromRequest(req fiber.Ctx) *Context {
	r := req.Request

	var ctx context.Context
	if r != nil {
		ctx = req.Context()
	}
	ctxi := reqctx.NewRequestContext[fiber.Ctx](ctx, req)
	setWithReq(ctxi, req.Request())
	return ctxi
}

func setWithReq(c *reqctx.Context[fiber.Ctx], r *fasthttp.Request) {
	if r == nil {
		return
	}
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = fasthttpctx.Device(&r.Header)
	c.Internal = stringsi.BytesToString(r.Header.Peek(httpi.HeaderGrpcInternal))
}

func DeviceFromHeader(r http.Header) *reqctx.DeviceInfo {
	return reqctx.Device(r.Get(httpi.HeaderDeviceInfo),
		r.Get(httpi.HeaderArea), r.Get(httpi.HeaderLocation),
		r.Get(httpi.HeaderUserAgent), r.Get(httpi.HeaderXForwardedFor))
}

type FiberContext Context

func (c *FiberContext) SetHeader(md metadata.MD) error {
	resp := c.RequestCtx.Response()
	for k, v := range md {
		if len(v) > 0 {
			resp.Header.Set(k, v[0])
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

func (c *FiberContext) SendHeader(md metadata.MD) error {
	resp := c.RequestCtx.Response()
	for k, v := range md {
		if len(v) > 0 {
			resp.Header.Set(k, v[0])
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

func (c *FiberContext) WriteHeader(k, v string) error {
	c.RequestCtx.Response().Header.Set(k, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FiberContext) SetCookie(v string) error {
	c.RequestCtx.Response().Header.Set(httpi.HeaderSetCookie, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FiberContext) SetTrailer(md metadata.MD) error {
	req := c.RequestCtx.Request()
	for k, v := range md {
		if len(v) > 0 {
			req.Header.Set(k, v[0])
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
