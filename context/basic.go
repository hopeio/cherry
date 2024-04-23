package context

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Context struct {
	ctx     context.Context
	TraceID string
	Values  map[string]any
}

func NewContext(ctx context.Context) *Context {
	var traceId string
	if ctx != nil {
		span := trace.SpanFromContext(ctx)
		if spanContext := span.SpanContext(); spanContext.IsValid() {
			traceId = spanContext.TraceID().String()
		}
	} else {
		ctx = context.Background()
	}

	if traceId == "" {
		traceId = uuid.New().String()
	}
	return &Context{
		ctx:     ctx,
		TraceID: traceId,
		Values:  map[string]any{},
	}
}

func (c *Context) StartSpan(name string, o ...trace.SpanStartOption) (*Context, trace.Span) {
	ctx, span := tracer.Start(c.ctx, name, o...)
	c.ctx = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID().String()
	}
	return c, span
}

func (c *Context) ContextWrapper() context.Context {
	return context.WithValue(c.ctx, ctxKey{}, c)
}

func ContextFromContext(ctx context.Context) *Context {
	if ctx == nil {
		return NewContext(context.Background())
	}

	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*Context)
	if !ok {
		c = NewContext(ctx)
	}
	c.ctx = ctx
	return c
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) SetContext(ctx context.Context) {
	c.ctx = ctx
}
