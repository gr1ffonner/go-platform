package clickhouse

import (
	"context"
	"fmt"
	"go-platform/internal/models/dogs"
	"go-platform/pkg/db/clickhouse"
	"go-platform/pkg/metrics"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type ClickHouseRepositoryMetricsInterface interface {
	RecordQuery(operation, table string, duration time.Duration)
	RecordError(operation, table, errorType string)
}

type ClickHouseRepository struct {
	clickhouse *clickhouse.ClickHouseClient
	dbMetrics  ClickHouseRepositoryMetricsInterface
}

func NewClickHouseRepository(clickhouse *clickhouse.ClickHouseClient, dbMetrics *metrics.DatabaseMetrics) *ClickHouseRepository {
	return &ClickHouseRepository{clickhouse: clickhouse, dbMetrics: dbMetrics}
}

// InsertDog inserts a dog into ClickHouse
func (r *ClickHouseRepository) InsertDog(ctx context.Context, dog *dogs.Dog) (string, error) {
	query := `
		INSERT INTO dogs (id, breed, image_url, created_at)
		VALUES (?, ?, ?, ?)`

	id := uuid.New().String()

	err := r.clickhouse.Conn().Exec(ctx, query, id, dog.Breed, dog.ImageURL, dog.CreatedAt)
	if err != nil {
		slog.Error("Failed to insert dog into ClickHouse", "error", err)
		return "", fmt.Errorf("failed to insert dog into ClickHouse: %w", err)
	}

	slog.Info("Dog inserted into ClickHouse", "id", id)
	return id, nil
}
