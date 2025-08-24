package main

import (
	"context"
	"errors"
	"fmt"
	restclientexample "go-platform/internal/clients/rest-client-example"
	"go-platform/internal/clients/s3"
	grpc "go-platform/internal/gprc"
	"go-platform/internal/handlers"
	"go-platform/internal/service/dogs"
	"go-platform/pkg/broker/nats"
	"go-platform/pkg/cache/redis"
	"go-platform/pkg/config"
	"go-platform/pkg/db/postgre"
	"go-platform/pkg/logger"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	_ "go-platform/api" // Import Swagger docs
)

// gracefulShutdown handles the graceful shutdown of all services
func gracefulShutdown(ctx context.Context, log *slog.Logger, services ...any) {
	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Info("Shutting down services...")

	// Shutdown each service
	for _, service := range services {
		switch s := service.(type) {
		case *http.Server:
			if err := s.Shutdown(shutdownCtx); err != nil {
				log.Error("HTTP server shutdown failed", "error", err)
			} else {
				log.Info("HTTP server stopped gracefully")
			}
		case interface{ GracefulStop() }:
			s.GracefulStop()
			log.Info("gRPC server stopped gracefully")
		case *postgre.PostgresClient:
			s.Close()
			log.Info("Database connection closed")
		case *redis.RedisClient:
			if err := s.Close(); err != nil {
				log.Error("Cache connection close failed", "error", err)
			} else {
				log.Info("Cache connection closed")
			}
		case *nats.NATSClient:
			s.Close()
			log.Info("Broker connection closed")
		}
	}

	log.Info("All services stopped")
}

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

	// Initialize unified logger
	logger.InitLogger(cfg.Logger)
	log := slog.Default()

	slog.Info("Config and logger initialized")

	// Create main context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database
	db, err := postgre.NewPostgres(ctx, cfg.Database.DSN)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	slog.Info("Database connected successfully")

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
	dogsService := dogs.NewDogsService(dogsAPI, s3Client)

	// Initialize handlers and router
	handler := handlers.NewHandler(dogsService)
	router := handlers.InitRouter(handler)

	grpcPort := cfg.Server.GRPCPort
	httpPort := cfg.Server.HTTPPort

	// HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: router,
	}

	// gRPC server
	grpcServer := grpc.NewServer()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Error("failed to create gRPC listener", "error", err)
		os.Exit(1)
	}

	// Setup signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server
	go func() {
		log.Info("Starting HTTP server", "port", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", "error", err)
			cancel()
		}
	}()

	// Start gRPC server
	go func() {
		log.Info("Starting gRPC server", "port", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Error("failed to serve gRPC", "error", err)
			cancel()
		}
	}()

	log.Info("All servers are ready to handle requests")

	// Wait for shutdown signal
	select {
	case <-stop:
		log.Info("Shutdown signal received")
	case <-ctx.Done():
		log.Info("Context canceled")
	}

	// Graceful shutdown
	gracefulShutdown(ctx, log, httpServer, grpcServer, db, cache, broker)
}
