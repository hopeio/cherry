package context

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

var (
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	meter          metric.Meter
	meterProvider  *sdkmetric.MeterProvider
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

func SetTracer(t trace.Tracer) {
	tracer = t
}

func SetMeter(m metric.Meter) {
	meter = m
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
