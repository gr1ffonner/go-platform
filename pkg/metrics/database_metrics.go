package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DatabaseMetrics holds database-related metrics
type DatabaseMetrics struct {
	DatabaseConnections   prometheus.Gauge
	DatabaseQueryDuration *prometheus.HistogramVec
	DatabaseErrorRate     *prometheus.CounterVec
}

// NewDatabaseMetrics creates a new database metrics instance
func NewDatabaseMetrics(registry *prometheus.Registry) *DatabaseMetrics {
	return &DatabaseMetrics{
		DatabaseConnections: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections_active",
				Help: "Number of active database connections",
			},
		),

		DatabaseQueryDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),

		DatabaseErrorRate: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_errors_total",
				Help: "Total number of database errors",
			},
			[]string{"operation", "table", "error_type"},
		),
	}
}

// RecordQuery records database query metrics
func (d *DatabaseMetrics) RecordQuery(operation, table string, duration time.Duration) {
	d.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordError records database error metrics
func (d *DatabaseMetrics) RecordError(operation, table, errorType string) {
	d.DatabaseErrorRate.WithLabelValues(operation, table, errorType).Inc()
}

// SetConnections sets the current number of database connections
func (d *DatabaseMetrics) SetConnections(count float64) {
	d.DatabaseConnections.Set(count)
}
