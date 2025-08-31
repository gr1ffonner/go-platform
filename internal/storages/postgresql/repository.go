package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go-platform/internal/models/dogs"
	"go-platform/pkg/db/postgre"
	"go-platform/pkg/metrics"
)

type PostgresRepositoryMetricsInterface interface {
	RecordQuery(operation, table string, duration time.Duration)
	RecordError(operation, table, errorType string)
}

type PostgresRepository struct {
	postgres  *postgre.PostgresClient
	dbMetrics PostgresRepositoryMetricsInterface
}

func NewPostgresRepository(postgres *postgre.PostgresClient, dbMetrics *metrics.DatabaseMetrics) *PostgresRepository {
	return &PostgresRepository{postgres: postgres, dbMetrics: dbMetrics}
}

// InsertDog inserts a dog into PostgreSQL
func (r *PostgresRepository) InsertDog(ctx context.Context, dog *dogs.Dog) (string, error) {
	start := time.Now()
	defer func() {
		r.dbMetrics.RecordQuery("insert", "dogs", time.Since(start))
	}()

	query := `
		INSERT INTO dogs (breed, image_url, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	dog.CreatedAt = time.Now()

	var id int
	err := r.postgres.Pool().QueryRow(ctx, query, dog.Breed, dog.ImageURL, dog.CreatedAt).Scan(&id)
	if err != nil {
		slog.Error("Failed to insert dog into PostgreSQL", "error", err)
		return "", fmt.Errorf("failed to insert dog into PostgreSQL: %w", err)
	}

	slog.Info("Dog inserted into PostgreSQL", "id", id)

	return fmt.Sprintf("%d", id), nil
}
