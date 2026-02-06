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
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
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
		if c.Propagator == nil {
			c.Propagator = c.newPropagator()
		}
		// Set up Propagator.
		otel.SetTextMapPropagator(c.Propagator)
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
		if c.TracerProvider == nil {
			c.TracerProvider, err = c.newTraceProvider(ctx, res)
			if err != nil {
				handleErr(err)
				return
			}
		}
		shutdownFuncs = append(shutdownFuncs, c.TracerProvider.Shutdown)
		otel.SetTracerProvider(c.TracerProvider)

		if c.MeterProvider == nil {
			// Set up meter provider.
			c.MeterProvider, err = c.newMeterProvider(ctx, res)
			if err != nil {
				handleErr(err)
				return
			}
		}
		shutdownFuncs = append(shutdownFuncs, c.MeterProvider.Shutdown)
		otel.SetMeterProvider(c.MeterProvider)
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
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter, c.BatchSpanProcessorOpts...),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
	), nil
}

func (c *TelemetryConfig) newMeterProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {

	options := []sdkmetric.Option{sdkmetric.WithResource(res)}
	if len(c.PrometheusExportOpts) > 0 {
		var err error
		reader, err := prometheus.New(c.PrometheusExportOpts...)
		if err != nil {
			return nil, err
		}
		options = append(options, sdkmetric.WithReader(reader))
	}
	if len(c.StdoutExportOpts) > 0 {
		exporter, err := stdoutmetric.New(c.StdoutExportOpts...)
		if err != nil {
			return nil, err
		}
		reader := sdkmetric.NewPeriodicReader(exporter, c.PeriodicReaderOps...)
		options = append(options, sdkmetric.WithReader(reader))
	}

	return sdkmetric.NewMeterProvider(options...), nil
}
