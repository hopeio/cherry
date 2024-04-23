package context

import (
	"context"
	"github.com/google/uuid"
	"github.com/hopeio/cherry/utils/net/http"
	timei "github.com/hopeio/cherry/utils/time"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"runtime"
	"strings"
	"sync"
	"time"
)

func GetPool[REQ any]() sync.Pool {
	return sync.Pool{New: func() any {
		return new(RequestContext[REQ])
	}}
}

type RequestContext[REQ any] struct {
	ctx context.Context

	TraceID string

	Token  string
	AuthID string
	AuthInfo
	AuthInfoRaw string

	*DeviceInfo

	http.RequestAt
	RequestCtx REQ
	grpc.ServerTransportStream

	Internal string
	Values   map[string]any
}

func (c *RequestContext[REQ]) StartSpan(name string, o ...trace.SpanStartOption) (*RequestContext[REQ], trace.Span) {
	ctx, span := tracer.Start(c.ctx, name, o...)
	c.ctx = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID().String()
	}
	return c, span
}

func Tracing(ctx context.Context, name string) (context.Context, trace.Span) {

	if span := trace.SpanFromContext(ctx); span != nil {
		return ctx, span
	}
	if name == "" {
		pc, _, _, _ := runtime.Caller(3)
		name = runtime.FuncForPC(pc).Name()
	}

	return tracer.Start(ctx, name)
}

func methodFamily(m string) string {
	m = strings.TrimPrefix(m, "/") // remove leading slash
	if i := strings.Index(m, "/"); i >= 0 {
		m = m[:i] // remove everything from second slash
	}
	return m
}

type ctxKey struct{}

func (c *RequestContext[REQ]) ContextWrapper() context.Context {
	return context.WithValue(c.ctx, ctxKey{}, c)
}

func RequestContextFromContext[REQ any](ctx context.Context) *RequestContext[REQ] {
	if ctx == nil {
		return NewRequestContext[REQ](context.Background(), "")
	}

	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*RequestContext[REQ])
	if !ok {
		var traceId string
		span := trace.SpanFromContext(ctx)
		if spanContext := span.SpanContext(); spanContext.IsValid() {
			traceId = spanContext.TraceID().String()
		}
		c = NewRequestContext[REQ](ctx, traceId)
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	c.ctx = ctx
	return c
}

func (c *RequestContext[REQ]) Context() context.Context {
	return c.ctx
}

func (c *RequestContext[REQ]) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func NewRequestContext[REQ any](ctx context.Context, traceId string) *RequestContext[REQ] {
	now := time.Now()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	return &RequestContext[REQ]{
		ctx:     ctx,
		TraceID: traceId,
		RequestAt: http.RequestAt{
			Time:       now,
			TimeStamp:  now.Unix(),
			TimeString: now.Format(timei.TimeFormat),
		},
		ServerTransportStream: grpc.ServerTransportStreamFromContext(ctx),
		Values:                map[string]any{},
	}
}

func (c *RequestContext[REQ]) reset(ctx context.Context) *RequestContext[REQ] {
	span := trace.SpanFromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID().String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	c.ctx = ctx
	c.RequestAt.Time = now
	c.RequestAt.TimeString = now.Format(timei.TimeFormat)
	c.RequestAt.TimeStamp = now.Unix()
	return c
}

func (c *RequestContext[REQ]) Method() string {
	if c.ServerTransportStream != nil {
		return c.ServerTransportStream.Method()
	}
	return ""
}
