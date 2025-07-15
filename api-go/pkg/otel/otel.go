package otel

import (
	"context"
	"errors"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// New creates a new OTel provider and registers it globally.
func New(ctx context.Context) (shutdown func(context.Context) error, err error) {
	// shutdown
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error

		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Initialize tracer
	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Initialize meter
	resource, _ := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME")),
			semconv.ServiceVersion("1.0"),
		))

	meterProvider, err := newMeterProvider(ctx, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Initialize custom metrics
	meter := meterProvider.Meter("api-go")
	if _, err = newMetrics(meter); err != nil {
		handleErr(err)
		return
	}

	return
}

// newTraceProvider creates a trace exporter that forward traces to OpenTelemetry Collector.
func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	// exporter
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	// provider
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(time.Second)),
	)
	return provider, nil
}

// newMeterProvider creates a meter exporter that outputs metrics to OpenTelemetry Collector.
func newMeterProvider(ctx context.Context, resource *resource.Resource) (*metric.MeterProvider, error) {
	// exporter
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	// provider
	return metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	), nil
}
