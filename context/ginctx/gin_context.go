package ginctx

import (
	"context"
	"github.com/gin-gonic/gin"
	contexti "github.com/hopeio/cherry/context"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc/metadata"
	"net/http"
)

type Context = contexti.RequestContext[*gin.Context]

func ContextFromContext(ctx context.Context) *Context {
	return contexti.RequestContextFromContext[*gin.Context](ctx)
}

func ContextFromRequest(req *gin.Context, tracing bool) (*Context, trace.Span) {
	r := req.Request

	ctx := context.Background()
	if r != nil {
		ctx = r.Context()
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
			/*	traceString := r.Header.Get(httpi.HeaderGrpcTraceBin)
				if traceString == "" {
					traceString = r.Header.Get(httpi.HeaderTraceBin)
				}
			*/
			ctx, span = contexti.Tracing(ctx, r.RequestURI)
		} else {
			ctx, span = contexti.Tracing(ctx, "")
		}
		if spanContext := span.SpanContext(); spanContext.IsValid() {
			traceId = spanContext.TraceID().String()
		}
	}

	ctxi := contexti.NewRequestContext[*gin.Context](ctx, traceId)
	ctxi.RequestCtx = req
	setWithHttpReq(ctxi, r)
	return ctxi, span
}

func setWithHttpReq(c *contexti.RequestContext[*gin.Context], r *http.Request) {
	if r == nil {
		return
	}
	c.DeviceInfo = DeviceFromHeader(r.Header)
	c.Internal = r.Header.Get(httpi.HeaderGrpcInternal)
	c.Token = httpi.GetToken(r)
}

func DeviceFromHeader(r http.Header) *contexti.DeviceInfo {
	return contexti.Device(r.Get(httpi.HeaderDeviceInfo),
		r.Get(httpi.HeaderArea), r.Get(httpi.HeaderLocation),
		r.Get(httpi.HeaderUserAgent), r.Get(httpi.HeaderXForwardedFor))
}

type GinContext contexti.RequestContext[*gin.Context]

func (c *GinContext) SetHeader(md metadata.MD) error {
	header := c.RequestCtx.Writer.Header()
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

func (c *GinContext) SendHeader(md metadata.MD) error {
	header := c.RequestCtx.Writer.Header()
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

func (c *GinContext) WriteHeader(k, v string) error {
	c.RequestCtx.Writer.Header().Set(k, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GinContext) SetCookie(v string) error {
	c.RequestCtx.Writer.Header().Set(httpi.HeaderSetCookie, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GinContext) SetTrailer(md metadata.MD) error {
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
