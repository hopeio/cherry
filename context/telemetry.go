package context

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

var (
	tracer trace.Tracer
	meter  metric.Meter
)

const (
	KindKey = attribute.Key("context.key")
)

var (
	serverKeyValue = KindKey.String("server")
)

func init() {
	tracer = otel.Tracer("context")
	meter = otel.Meter("context")
}

func Tracer() trace.Tracer {
	return tracer
}

func SetTracer(t trace.Tracer) {
	tracer = t
}

func Meter() metric.Meter {
	return meter
}

func SetMeter(m metric.Meter) {
	meter = m
}

func StartSpan(ctx context.Context, name string, o ...trace.SpanStartOption) (context.Context, trace.Span) {
	if name == "" {
		pc, _, _, _ := runtime.Caller(3)
		name = runtime.FuncForPC(pc).Name()
	}
	return tracer.Start(ctx, name, o...)
}

func (c *Context) StartSpan(name string, o ...trace.SpanStartOption) trace.Span {
	ctx, span := tracer.Start(c.ctx, name, o...)
	c.ctx = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID().String()
	}
	return span
}

func (c *Context) StartSpanX(name string, o ...trace.SpanStartOption) (*Context, trace.Span) {
	ctx, span := tracer.Start(c.ctx, name, o...)
	c.ctx = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID().String()
	}
	return c, span
}

func (c *Context) StartSpanEnd(name string, o ...trace.SpanStartOption) func(options ...trace.SpanEndOption) {
	ctx, span := tracer.Start(c.ctx, name, o...)
	c.ctx = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID().String()
	}
	return span.End
}
