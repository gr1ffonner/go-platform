package main

import (
	"context"
	"go-platform/pkg/broker/nats"
	"go-platform/pkg/cache/redis"
	"go-platform/pkg/config"
	"go-platform/pkg/db/postgre"
	"go-platform/pkg/logger"
	"log/slog"
)

// @title	    Go Platform
// @version		1.0
// @description	Go Platform API
func main() {
	// Load config first
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	// Initialize unified logger
	logger.InitLogger(cfg.Logger)
	log := slog.Default()

	slog.Info("Config and logger initialized")

	ctx := context.Background()

	// Initialize database
	db, err := postgre.NewPostgres(ctx, cfg.Database.DSN)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer db.Close()
	log.Info("Database connected successfully")

	// Initialize cache
	cache, err := redis.NewRedis(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}
	defer cache.Close()
	log.Info("Redis connected successfully")

	// Initialize broker
	broker, err := nats.NewNATS(ctx, cfg.NATS.URL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		panic(err)
	}
	defer broker.Close()
	log.Info("NATS connected successfully")
}
