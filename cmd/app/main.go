package main

import (
	"context"
	restclientexample "go-platform/internal/clients/rest-client-example"
	"go-platform/internal/clients/s3"
	grpc "go-platform/internal/gprc"
	"go-platform/internal/handlers"
	"go-platform/internal/services/dogs"
	"go-platform/pkg/broker/nats"
	"go-platform/pkg/cache/redis"
	"go-platform/pkg/config"
	"go-platform/pkg/logger"
	"go-platform/pkg/server"
	"go-platform/pkg/utils"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "go-platform/api" // Import Swagger docs
)

// @title			Go Platform
// @version		1.0
// @description	Go Platform API
func main() {
	// Load config first
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}
	slog.Info("Config loaded", "config", cfg)

	// Initialize unified logger
	logger.InitLogger(cfg.Logger)
	log := slog.Default()

	slog.Info("Config and logger initialized")

	// Create main context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Storage layer initializing
	storage, err := utils.GetStorage(ctx, cfg)
	if err != nil {
		slog.Error("Failed to get storage", "error", err)
		panic(err)
	}

	// Initialize S3 client
	s3Client, err := s3.NewClientS3(
		cfg.S3.KeyID,
		cfg.S3.KeySecret,
		cfg.S3.Bucket,
		cfg.S3.BaseEndpoint,
		cfg.S3.BasePublicEndpoint,
		cfg.S3.Region,
	)
	if err != nil {
		log.Error("Failed to connect to S3", "error", err)
		panic(err)
	}

	slog.Info("S3 client connected successfully")

	// Initialize cache
	cache, err := redis.NewRedis(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}

	slog.Info("Cache connected successfully")

	// Initialize broker
	broker, err := nats.NewNATS(ctx, cfg.NATS.URL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		panic(err)
	}

	slog.Info("Broker connected successfully")

	// Initialize dog API client
	dogsAPI := restclientexample.NewDogAPI()

	// Initialize dogs service
	dogsService := dogs.NewDogsService(dogsAPI, s3Client, storage.Repository)

	// Initialize handlers
	handler := handlers.NewHandler(dogsService)

	// gRPC server
	grpcServer := grpc.NewServer(dogsService)

	// Create unified server first to get metrics
	srv, err := server.NewServer(cfg, nil, grpcServer)
	if err != nil {
		log.Error("Failed to create server", "error", err)
		panic(err)
	}

	// Initialize router with metrics
	router := handlers.InitRouter(handler, srv.Metrics.HTTP)

	// Update server with the router
	srv.HTTP.Handler = router

	// Setup signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server
	if err := srv.Start(ctx); err != nil {
		log.Error("Failed to start server", "error", err)
		panic(err)
	}

	log.Info("All servers are ready to handle requests")

	// Wait for shutdown signal
	select {
	case <-stop:
		log.Info("Shutdown signal received")
	case <-ctx.Done():
		log.Info("Context canceled")
	}

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "error", err)
	}

	// Additional cleanup
	utils.GracefulShutdown(ctx, nil, nil, storage.DBClient, cache, broker)
}
