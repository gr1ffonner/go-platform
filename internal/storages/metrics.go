package storages

import "time"

// DatabaseMetricsInterface defines the interface for database metrics collection
type DatabaseMetricsInterface interface {
	// Basic metrics
	RecordQuery(operation, table string, duration time.Duration)
	RecordError(operation, table, errorType string)

	// Extended metrics with additional labels
	RecordQueryWithBreed(operation, table, breed string, duration time.Duration)
	RecordErrorWithBreed(operation, table, breed, errorType string)

	// Generic metrics with custom labels
	RecordQueryWithLabels(operation, table string, labels map[string]string, duration time.Duration)
	RecordErrorWithLabels(operation, table, errorType string, labels map[string]string)
}
