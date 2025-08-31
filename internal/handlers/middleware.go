package handlers

import (
	"log/slog"
	"net/http"

	"go-platform/pkg/metrics"
)

// LoggingMiddleware logs the details of each request and response
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "path", r.URL.Path)

		lrw := &LoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		slog.Info("Response status", "status", lrw.statusCode)
	})
}

// MetricsMiddleware creates HTTP middleware for metrics collection
func MetricsMiddleware(httpMetrics *metrics.HTTPMetrics) func(http.Handler) http.Handler {
	return httpMetrics.HTTPMiddleware
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
