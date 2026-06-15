// Package tracing configures OpenTelemetry distributed tracing. When an OTLP
// endpoint is configured it exports spans to that collector; otherwise it falls
// back to a stdout exporter (dev) or a no-op (when stdout export is disabled).
package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// Setup initialises the global OpenTelemetry tracer provider.
//
//   - endpoint = ""         → no-op provider (tests, no collector)
//   - endpoint = "stdout"   → stdout exporter (dev, human-readable)
//   - endpoint = "http://…" → OTLP/HTTP exporter (staging / production)
//
// Returns a shutdown function; defer it in main.
func Setup(ctx context.Context, serviceName, endpoint string) (shutdown func(context.Context) error, err error) {
	if endpoint == "" {
		// No-op: set a noop tracer so otel.Tracer() never panics.
		otel.SetTracerProvider(noop.NewTracerProvider())
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing: build resource: %w", err)
	}

	var exp sdktrace.SpanExporter
	if endpoint == "stdout" {
		exp, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else {
		exp, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpointURL(endpoint),
			otlptracehttp.WithInsecure(),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("tracing: init exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

// Tracer returns a named tracer from the global provider. Use throughout
// usecases and repos to add child spans to the incoming request trace.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
