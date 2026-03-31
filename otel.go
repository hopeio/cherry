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
	"go.opentelemetry.io/otel/sdk/trace"
)

const ScopeName = "github.com/hopeio/cherry"

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

	newPropagator()


	if otel.GetTracerProvider() == nil {
		var tracerProvider *trace.TracerProvider
		tracerProvider, err = newTraceProvider()
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	}

	if otel.GetMeterProvider() == nil {
		var meterProvider *metric.MeterProvider
		meterProvider, err = newMeterProvider()
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)

	}
	if len(shutdownFuncs) == 0 {
		shutdown = nil
	}
	return
}

func newPropagator() {
	if otel.GetTextMapPropagator() == nil {
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}
}

func newTraceProvider() (*trace.TracerProvider, error) {
	tracerProvider := trace.NewTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	meterProvider := metric.NewMeterProvider()
	otel.SetMeterProvider(meterProvider)
	return meterProvider, nil
}
