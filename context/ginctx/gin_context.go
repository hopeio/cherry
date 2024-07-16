package ginctx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/context/reqctx"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type Context = reqctx.Context[*gin.Context]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[*gin.Context](ctx)
}

func FromRequest(req *gin.Context) *Context {
	r := req.Request

	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	}

	ctxi := reqctx.NewRequestContext[*gin.Context](ctx, req)
	setWithHttpReq(ctxi, r)
	return ctxi
}

func setWithHttpReq(c *reqctx.Context[*gin.Context], r *http.Request) {
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

type GinContext Context

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
