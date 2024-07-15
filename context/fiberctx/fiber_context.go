package fiberctx

import (
	"context"
	"github.com/gofiber/fiber/v3"
	contexti "github.com/hopeio/cherry/context"
	"github.com/hopeio/cherry/context/fasthttpctx"
	httpi "github.com/hopeio/cherry/utils/net/http"
	fasthttpi "github.com/hopeio/cherry/utils/net/http/fasthttp"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc/metadata"
	"net/http"
)

type Context = contexti.RequestContext[fiber.Ctx]

func ContextFromContext(ctx context.Context) *Context {
	return contexti.RequestContextFromContext[fiber.Ctx](ctx)
}

func ContextFromRequest(req fiber.Ctx, tracing bool) (*Context, trace.Span) {
	r := req.Request

	ctx := context.Background()
	if r != nil {
		ctx = req.Context()
	}
	var traceId string
	var span trace.Span
	if tracing {

		if r != nil {
			// go.opencensus.io/trace 完全包含了golang.org/x/net/trace 的功能
			// grpc内置配合,看了源码并没有启用，根本没调用
			// 系统trace只能追踪单个请求，且只记录时间及是否完成，只能/debug/requests看
			/*			t = gtrace.New(methodFamily(r.RequestURI), r.RequestURI)
						ctx = gtrace.NewContext(ctx, t)
			*/

			// 直接从远程读取Trace信息，Trace是否为空交给propagation包判断
			/*	traceString := req.Get(httpi.HeaderGrpcTraceBin)
				if traceString == "" {
					traceString = req.Get(httpi.HeaderTraceBin)
				}*/

			ctx, span = contexti.Tracing(ctx, req.BaseURL())
		} else {
			ctx, span = contexti.Tracing(ctx, "")
		}
		if spanContext := span.SpanContext(); spanContext.IsValid() {
			traceId = spanContext.TraceID().String()
		}
	}

	ctxi := contexti.NewRequestContext[fiber.Ctx](ctx, req, traceId)
	setWithReq(ctxi, req.Request())
	return ctxi, span
}

func setWithReq(c *contexti.RequestContext[fiber.Ctx], r *fasthttp.Request) {
	if r == nil {
		return
	}
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = fasthttpctx.Device(&r.Header)
	c.Internal = stringsi.BytesToString(r.Header.Peek(httpi.HeaderGrpcInternal))
}

func DeviceFromHeader(r http.Header) *contexti.DeviceInfo {
	return contexti.Device(r.Get(httpi.HeaderDeviceInfo),
		r.Get(httpi.HeaderArea), r.Get(httpi.HeaderLocation),
		r.Get(httpi.HeaderUserAgent), r.Get(httpi.HeaderXForwardedFor))
}

type FiberContext contexti.RequestContext[fiber.Ctx]

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
