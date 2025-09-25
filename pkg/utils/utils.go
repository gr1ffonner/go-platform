package utils

import (
	"context"
	"fmt"
	"go-platform/internal/models/dogs"
	clickhouseRepo "go-platform/internal/storages/clickhouse"
	"go-platform/internal/storages/postgresql"
	"go-platform/pkg/config"
	"go-platform/pkg/db/clickhouse"
	"go-platform/pkg/db/postgre"
	"go-platform/pkg/metrics"
	"log/slog"
	"time"
)

type RepositoryMetricsInterface interface {
	RecordQuery(operation, table string, duration time.Duration)
	RecordError(operation, table, errorType string)
}

type Repository interface {
	// return string due to clickhouse dont have auto increment and
	// we should use uuid for simple row
	InsertDog(ctx context.Context, dog *dogs.Dog) (string, error)
}

// StorageResult contains both repository and database client for proper cleanup
type Storage struct {
	Repository Repository
	DBClient   interface{}
}

func GetStorage(ctx context.Context, cfg *config.Config, dbMetrics *metrics.DatabaseMetrics) (*Storage, error) {
	// Initialize storage
	switch cfg.Server.Storage {
	case "postgres":
		if cfg.Database.PostgresDSN == "" {
			return nil, fmt.Errorf("POSTGRES_DSN is required for postgres storage")
		}
		pgStorage, err := postgre.NewPostgres(ctx, cfg.Database.PostgresDSN)
		if err != nil {
			slog.Error("Failed to connect to postgres", "error", err)
			return nil, err
		}
		slog.Info("PostgreSQL connected successfully")

		pgRepository := postgresql.NewPostgresRepository(pgStorage, dbMetrics)
		slog.Info("PostgreSQL repository initialized")

		return &Storage{
			Repository: pgRepository,
			DBClient:   pgStorage,
		}, nil

	case "clickhouse":
		if cfg.Database.ClickHouseDSN == "" {
			return nil, fmt.Errorf("CLICKHOUSE_DSN is required for clickhouse storage")
		}
		clickhouseStorage, err := clickhouse.NewClickHouse(ctx, cfg.Database.ClickHouseDSN)
		if err != nil {
			slog.Error("Failed to connect to clickhouse", "error", err)
			return nil, err
		}
		slog.Info("ClickHouse connected successfully")

		clickhouseRepository := clickhouseRepo.NewClickHouseRepository(clickhouseStorage, dbMetrics)
		slog.Info("ClickHouse repository initialized")

		return &Storage{
			Repository: clickhouseRepository,
			DBClient:   clickhouseStorage,
		}, nil

	default:
		slog.Error("Invalid storage type", "storage", cfg.Server.Storage)
		return nil, fmt.Errorf("invalid storage type: %s, supported types: postgres, mysql, clickhouse", cfg.Server.Storage)
	}
}
