package tracer

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"go-platform/pkg/config"
)

// Tracer wraps the OpenTelemetry tracer
type Tracer struct {
	tracerProvider *sdktrace.TracerProvider
	shutdown       func(context.Context) error
}

// NewTracer creates a new OpenTelemetry tracer
func NewTracer(ctx context.Context, cfg config.MetricsProviderConfig) (*Tracer, error) {
	// Create OTLP exporter
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint))
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1.0)), // 100% sampling for development
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// Set global text map propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	slog.Info("OpenTelemetry tracer initialized",
		"service", cfg.ServiceName,
		"version", cfg.ServiceVersion,
		"endpoint", cfg.OTLPEndpoint,
	)

	return &Tracer{
		tracerProvider: tracerProvider,
		shutdown:       tracerProvider.Shutdown,
	}, nil
}

// Shutdown gracefully shuts down the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.shutdown != nil {
		return t.shutdown(ctx)
	}
	return nil
}

// GetTracer returns the global tracer instance
func GetTracer() trace.Tracer {
	return otel.Tracer("go-platform")
}

// StartSpan creates a new span with the given name and options
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := GetTracer()
	return tracer.Start(ctx, name, opts...)
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanError marks the span as failed and adds error information
func SetSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
