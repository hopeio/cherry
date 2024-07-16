package reqctx

import (
	"context"
	"github.com/google/uuid"
	context2 "github.com/hopeio/cherry/context"
	"github.com/hopeio/cherry/utils/net/http"
	timei "github.com/hopeio/cherry/utils/time"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"strings"
	"sync"
	"time"
)

func GetPool[REQ any]() sync.Pool {
	return sync.Pool{New: func() any {
		return new(Context[REQ])
	}}
}

type ReqValue[REQ any] struct {
	Token       string
	AuthInfoRaw string
	AuthID      string
	AuthInfo
	*DeviceInfo
	RequestCtx REQ
	grpc.ServerTransportStream
	Internal string
	http.RequestAt
}

type ReqCtx[REQ any] context2.ValueContext[ReqValue[REQ]]

type Context[REQ any] struct {
	context2.Context
	ReqValue[REQ]
}

func methodFamily(m string) string {
	m = strings.TrimPrefix(m, "/") // remove leading slash
	if i := strings.Index(m, "/"); i >= 0 {
		m = m[:i] // remove everything from second slash
	}
	return m
}

func (c *Context[REQ]) Wrapper() context.Context {
	return context.WithValue(c.Context.BaseContext(), context2.WrapperKey(), c)
}

func (c *Context[REQ]) StartSpanX(name string, o ...trace.SpanStartOption) (*Context[REQ], trace.Span) {
	span := c.Context.StartSpan(name, o...)
	return c, span
}

func FromContextValue[REQ any](ctx context.Context) *Context[REQ] {
	if ctx == nil {
		return NewRequestContext[REQ](context.Background(), *new(REQ))
	}

	ctxi := ctx.Value(context2.WrapperKey())
	c, ok := ctxi.(*Context[REQ])
	if !ok {
		c = NewRequestContext[REQ](ctx, *new(REQ))
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	c.SetBaseContext(ctx)
	return c
}

func NewRequestContext[REQ any](ctx context.Context, req REQ) *Context[REQ] {
	now := time.Now()

	return &Context[REQ]{
		Context: *context2.NewContext(ctx),
		ReqValue: ReqValue[REQ]{
			RequestCtx: req,
			RequestAt: http.RequestAt{
				Time:       now,
				TimeStamp:  now.Unix(),
				TimeString: now.Format(timei.LayoutTimeMacro),
			},
			ServerTransportStream: grpc.ServerTransportStreamFromContext(ctx),
		},
	}
}

func (c *Context[REQ]) reset(ctx context.Context) *Context[REQ] {
	span := trace.SpanFromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID().String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	c.SetBaseContext(ctx)
	c.RequestAt.Time = now
	c.RequestAt.TimeString = now.Format(timei.LayoutTimeMacro)
	c.RequestAt.TimeStamp = now.Unix()
	return c
}

func (c *Context[REQ]) Method() string {
	if c.ServerTransportStream != nil {
		return c.ServerTransportStream.Method()
	}
	return ""
}
