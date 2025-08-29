package utils

import (
	"context"
	"fmt"
	"go-platform/internal/models/dogs"
	clickhouseRepo "go-platform/internal/storages/clickhouse"
	mysqlRepo "go-platform/internal/storages/mysql"
	"go-platform/internal/storages/postgresql"
	"go-platform/pkg/broker/nats"
	"go-platform/pkg/cache/redis"
	"go-platform/pkg/config"
	"go-platform/pkg/db/clickhouse"
	"go-platform/pkg/db/mysql"
	"go-platform/pkg/db/postgre"
	"log/slog"
	"net/http"
	"time"
)

type Repository interface {
	// return string due to clickhouse dont have auto increment and
	// we should use uuid for simple row
	InsertDog(ctx context.Context, dog *dogs.Dog) (string, error)
}

// gracefulShutdown handles the graceful shutdown of all services
func GracefulShutdown(ctx context.Context, services ...any) {
	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	slog.Info("Shutting down services...")

	// Shutdown each service
	for _, service := range services {
		switch s := service.(type) {
		case *http.Server:
			if err := s.Shutdown(shutdownCtx); err != nil {
				slog.Error("HTTP server shutdown failed", "error", err)
			} else {
				slog.Info("HTTP server stopped gracefully")
			}
		case interface{ GracefulStop() }:
			s.GracefulStop()
			slog.Info("gRPC server stopped gracefully")
		case *postgre.PostgresClient:
			s.Close()
			slog.Info("Postgres connection closed")
		case *mysql.MySQLClient:
			s.Close()
			slog.Info("Mysql connection closed")
		case *clickhouse.ClickHouseClient:
			s.Close()
			slog.Info("Clickhouse connection closed")
		case *redis.RedisClient:
			if err := s.Close(); err != nil {
				slog.Error("Cache connection close failed", "error", err)
			} else {
				slog.Info("Cache connection closed")
			}
		case *nats.NATSClient:
			s.Close()
			slog.Info("Broker connection closed")
		}
	}

}

// StorageResult contains both repository and database client for proper cleanup
type Storage struct {
	Repository Repository
	DBClient   interface{}
}

func GetStorage(ctx context.Context, cfg *config.Config) (*Storage, error) {
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

		pgRepository := postgresql.NewPostgresRepository(pgStorage)
		slog.Info("PostgreSQL repository initialized")

		return &Storage{
			Repository: pgRepository,
			DBClient:   pgStorage,
		}, nil

	case "mysql":
		if cfg.Database.MySQLDSN == "" {
			return nil, fmt.Errorf("MYSQL_DSN is required for mysql storage")
		}
		mysqlStorage, err := mysql.NewMySQL(ctx, cfg.Database.MySQLDSN)
		if err != nil {
			slog.Error("Failed to connect to mysql", "error", err)
			return nil, err
		}
		slog.Info("MySQL connected successfully")

		mysqlRepository := mysqlRepo.NewMySQLRepository(mysqlStorage)
		slog.Info("MySQL repository initialized")

		return &Storage{
			Repository: mysqlRepository,
			DBClient:   mysqlStorage,
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

		clickhouseRepository := clickhouseRepo.NewClickHouseRepository(clickhouseStorage)
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
