package handlers

import (
	"log/slog"
	"net/http"
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

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}
