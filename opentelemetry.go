/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"errors"

	_ "github.com/hopeio/gox/net/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
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

	propagator := newPropagator()
	// Set up Propagator.
	otel.SetTextMapPropagator(propagator)
	var res *resource.Resource
	res, err = resource.New(
		ctx, resource.WithFromEnv(), // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),      // Discover and provide process information.
		resource.WithOS(),           // Discover and provide OS information.
		resource.WithContainer(),    // Discover and provide container information.
		resource.WithHost())         // Discover and provide host information.

	if err != nil {
		return nil, err
	}

	if otel.GetTracerProvider() == nil {
		var tracerProvider *trace.TracerProvider
		tracerProvider, err = newTraceProvider(res)
		if err != nil {
			handleErr(err)
			return
		}

		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	if otel.GetMeterProvider() == nil {
		var meterProvider *metric.MeterProvider
		meterProvider, err = newMeterProvider(res)
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}
	if len(shutdownFuncs) == 0 {
		shutdown = nil
	}
	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	return trace.NewTracerProvider(trace.WithResource(res)), nil
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {

	return metric.NewMeterProvider(metric.WithResource(res)), nil
}
