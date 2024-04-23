package context

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
)

const (
	KindKey = attribute.Key("context.key")
)

var (
	serverKeyValue = KindKey.String("server")
)

func init() {
	tracer = otel.Tracer("context")
}

func SetTracer(t trace.Tracer) {
	tracer = t
}
