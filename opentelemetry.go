/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func (c *TelemetryConfig) setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
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
	if c.Enabled {
		if c.propagator == nil {
			c.propagator = c.newPropagator()
		}
		// Set up propagator.
		otel.SetTextMapPropagator(c.propagator)
		var res *resource.Resource
		res, err = resource.New(ctx) //resource.WithFromEnv(), // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		//resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		//resource.WithProcess(),      // Discover and provide process information.
		//resource.WithOS(),           // Discover and provide OS information.
		//resource.WithContainer(), // Discover and provide container information.
		//resource.WithHost(),         // Discover and provide host information.

		if err != nil {
			return nil, err
		}
		if c.tracerProvider == nil {
			c.tracerProvider, err = c.newTraceProvider(ctx, res)
			if err != nil {
				handleErr(err)
				return
			}
		}
		shutdownFuncs = append(shutdownFuncs, c.tracerProvider.Shutdown)
		otel.SetTracerProvider(c.tracerProvider)

		if c.meterProvider == nil {
			// Set up meter provider.
			c.meterProvider, err = c.newMeterProvider(ctx, res)
			if err != nil {
				handleErr(err)
				return
			}
		}
		shutdownFuncs = append(shutdownFuncs, c.meterProvider.Shutdown)
		otel.SetMeterProvider(c.meterProvider)
	}
	return
}

func (c *TelemetryConfig) newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (c *TelemetryConfig) newTraceProvider(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
	//stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
	), nil
}

func (c *TelemetryConfig) newMeterProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {

	var reader sdkmetric.Reader
	if c.EnablePrometheus {
		var err error
		reader, err = prometheus.New(c.prometheusOpts...)
		if err != nil {
			return nil, err
		}
	} else {
		exporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		reader = sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(time.Minute))
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
