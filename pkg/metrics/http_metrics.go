package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPMetrics holds HTTP-related metrics
type HTTPMetrics struct {
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge
	HTTPErrorRate        *prometheus.CounterVec
}

// NewHTTPMetrics creates a new HTTP metrics instance
func NewHTTPMetrics(registry *prometheus.Registry) *HTTPMetrics {
	return &HTTPMetrics{
		HTTPRequestsTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status_code"},
		),

		HTTPRequestDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status_code"},
		),

		HTTPRequestsInFlight: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),

		HTTPErrorRate: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors (4xx, 5xx)",
			},
			[]string{"method", "path", "status_code"},
		),
	}
}

// HTTPMiddleware creates HTTP middleware for metrics collection
func (h *HTTPMetrics) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment in-flight requests
		h.HTTPRequestsInFlight.Inc()
		defer h.HTTPRequestsInFlight.Dec()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		statusCode := fmt.Sprintf("%d", wrapped.statusCode)

		// Track all requests
		h.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
		h.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(duration)

		// Track errors (4xx, 5xx)
		if wrapped.statusCode >= 400 {
			h.HTTPErrorRate.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
