package server

import (
	"context"
	"errors"
	"github.com/hopeio/cherry/utils/log"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	meter      = otel.Meter("service")
	apiCounter metric.Int64Counter
	histogram  metric.Float64Histogram
)

func init() {
	var err error
	apiCounter, err = meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		log.Fatal(err)
	}
	histogram, err = meter.Float64Histogram(
		"task.duration",
		metric.WithDescription("The duration of task execution."),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context, enablePrometheus bool) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	otel.SetTextMapPropagator(newPropagator())

	// Set up trace provider.
	tracerProvider, err1 := newTraceProvider(ctx)
	if err != nil {
		handleErr(err1)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err1 := newMeterProvider(ctx, enablePrometheus)
	if err != nil {
		handleErr(err1)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
	//stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx) //resource.WithFromEnv(), // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
	//resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
	//resource.WithProcess(),      // Discover and provide process information.
	//resource.WithOS(),           // Discover and provide OS information.
	//resource.WithContainer(), // Discover and provide container information.
	//resource.WithHost(),         // Discover and provide host information.

	if err != nil {
		return nil, err
	}
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.5)),
	), nil
}

func newMeterProvider(ctx context.Context, enablePrometheus bool) (*sdkmetric.MeterProvider, error) {
	res, err := resource.New(ctx)
	var reader sdkmetric.Reader
	if enablePrometheus {
		reader, err = prometheus.New()
		if err != nil {
			return nil, err
		}
	} else {
		exporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		reader = sdkmetric.NewPeriodicReader(exporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			sdkmetric.WithInterval(3*time.Second))
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
		/*		sdkmetric.WithView(sdkmetric.NewView(
				sdkmetric.Instrument{Name: "histogram_*"},
				sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
					Boundaries: []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
				}},
			)),*/
	), nil
}
