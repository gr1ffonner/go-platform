package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"go-platform/pkg/config"
)

// Metrics holds all the application metrics
type Metrics struct {
	HTTP     *HTTPMetrics
	Database *DatabaseMetrics
	System   *SystemMetrics

	// Prometheus registry
	registry *prometheus.Registry
}

func NewMetrics(cfg config.MetricsProviderConfig) (*Metrics, error) {
	// Create Prometheus registry
	registry := prometheus.NewRegistry()

	// Create metrics
	metrics := &Metrics{
		HTTP:     NewHTTPMetrics(registry),
		Database: NewDatabaseMetrics(registry),
		System:   NewSystemMetrics(registry),
		registry: registry,
	}

	return metrics, nil
}

// NewOTLPMetrics creates OpenTelemetry metrics exporter
func NewOTLPMetrics(ctx context.Context, cfg config.MetricsProviderConfig) (*metric.MeterProvider, error) {
	// Create OTLP exporter
	var opts []otlpmetricgrpc.Option
	opts = append(opts, otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint))
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}
	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metrics exporter: %w", err)
	}

	// Create resource
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

	// Create meter provider
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	slog.Info("OpenTelemetry metrics initialized", "endpoint", cfg.OTLPEndpoint)
	return meterProvider, nil
}

// StartPrometheusServer starts the Prometheus metrics server
func (m *Metrics) StartPrometheusServer(ctx context.Context, port string) error {
	// Create HTTP handler for Prometheus metrics
	handler := promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Starting Prometheus metrics server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Prometheus metrics server error", "error", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
