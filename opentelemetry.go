/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"errors"
	"github.com/hopeio/utils/log"
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
func setupOTelSDK(ctx context.Context, config *TelemetryConfig) (shutdown func(context.Context) error, err error) {
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
	if config.EnableTracing {
		if config.propagator == nil {
			config.propagator = newPropagator()
		}
		// Set up propagator.
		otel.SetTextMapPropagator(config.propagator)

		if config.tracerProvider == nil {
			config.tracerProvider, err = newTraceProvider(ctx)
			if err != nil {
				handleErr(err)
				return
			}
		}
		shutdownFuncs = append(shutdownFuncs, config.tracerProvider.Shutdown)
		otel.SetTracerProvider(config.tracerProvider)
	}
	if config.EnableMetrics {
		if config.meterProvider == nil {
			// Set up meter provider.
			config.meterProvider, err = newMeterProvider(ctx)
			if err != nil {
				handleErr(err)
				return
			}
		}
		shutdownFuncs = append(shutdownFuncs, config.meterProvider.Shutdown)
		otel.SetMeterProvider(config.meterProvider)
	}
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

func newMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	res, err := resource.New(ctx)
	var reader sdkmetric.Reader
	exporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}
	reader = sdkmetric.NewPeriodicReader(exporter,
		sdkmetric.WithInterval(time.Minute))

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
