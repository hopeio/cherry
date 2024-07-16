package context

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Context struct {
	ctx      context.Context
	rootSpan trace.Span
	TraceID  string
}

func NewContext(ctx context.Context) *Context {
	var traceId string
	var rootSpan trace.Span
	if ctx != nil {
		span := trace.SpanFromContext(ctx)
		if spanContext := span.SpanContext(); spanContext.IsValid() {
			traceId = spanContext.TraceID().String()
			rootSpan = span
		}
	} else {
		ctx = context.Background()
	}
	if traceId == "" || rootSpan == nil {
		ctx, rootSpan = StartSpan(ctx, "")
		traceId = rootSpan.SpanContext().TraceID().String()
	}

	return &Context{
		ctx:      ctx,
		rootSpan: rootSpan,
		TraceID:  traceId,
	}
}

type ctxKey struct{}

func (c *Context) Wrapper() context.Context {
	return context.WithValue(c.ctx, ctxKey{}, c)
}

func WrapperKey() ctxKey {
	return ctxKey{}
}

func FromContextValue(ctx context.Context) *Context {
	if ctx == nil {
		return NewContext(nil)
	}

	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*Context)
	if !ok {
		c = NewContext(ctx)
	}
	c.ctx = ctx
	return c
}

func (c *Context) BaseContext() context.Context {
	return c.ctx
}

func (c *Context) SetBaseContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Context) RootSpan() trace.Span {
	return c.rootSpan
}

type ValueContext[V any] struct {
	Context
	Value V
}
