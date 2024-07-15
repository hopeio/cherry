package fasthttpctx

import (
	"context"
	contexti "github.com/hopeio/cherry/context"
	httpi "github.com/hopeio/cherry/utils/net/http"
	fasthttpi "github.com/hopeio/cherry/utils/net/http/fasthttp"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc/metadata"
)

type Context = contexti.RequestContext[*fasthttp.RequestCtx]

func ContextFromRequest(req *fasthttp.RequestCtx, tracing bool) (*Context, trace.Span) {
	r := &req.Request

	ctx := context.Background()
	if r != nil {
		ctx = req
	}
	var traceId string
	var span trace.Span
	if tracing {

		if req != nil {
			// go.opencensus.io/trace 完全包含了golang.org/x/net/trace 的功能
			// grpc内置配合,看了源码并没有启用，根本没调用
			// 系统trace只能追踪单个请求，且只记录时间及是否完成，只能/debug/requests看
			/*			t = gtrace.New(methodFamily(r.RequestURI), r.RequestURI)
						ctx = gtrace.NewContext(ctx, t)
			*/

			// 直接从远程读取Trace信息，Trace是否为空交给propagation包判断
			/*		traceString := r.Header.Peek(httpi.HeaderGrpcTraceBin)
					if traceString == nil {
						traceString = r.Header.Peek(httpi.HeaderTraceBin)
					}*/

			ctx, span = contexti.Tracing(ctx, string(r.RequestURI()))
		} else {
			ctx, span = contexti.Tracing(ctx, "")
		}
		if spanContext := span.SpanContext(); spanContext.IsValid() {
			traceId = spanContext.TraceID().String()
		}
	}

	ctxi := contexti.NewRequestContext[*fasthttp.RequestCtx](ctx, req, traceId)
	setWithReq(ctxi, r)
	return ctxi, span
}

func setWithReq(c *Context, r *fasthttp.Request) {
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = Device(&r.Header)
	c.Internal = stringsi.BytesToString(r.Header.Peek(httpi.HeaderGrpcInternal))
}

func Device(r *fasthttp.RequestHeader) *contexti.DeviceInfo {
	return contexti.Device(stringsi.BytesToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.BytesToString(r.Peek(httpi.HeaderArea)),
		stringsi.BytesToString(r.Peek(httpi.HeaderLocation)),
		stringsi.BytesToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.BytesToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}

type FastHttpContext contexti.RequestContext[*fasthttp.RequestCtx]

func (c *FastHttpContext) SetHeader(md metadata.MD) error {
	header := c.RequestCtx.Response.Header
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
	header := c.RequestCtx.Response.Header
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
