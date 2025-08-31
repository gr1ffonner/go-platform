package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"go-platform/pkg/config"
	"go-platform/pkg/metrics"
	"go-platform/pkg/tracer"
)

// GRPCServer interface for gRPC server operations
type GRPCServer interface {
	Serve(net.Listener) error
	GracefulStop()
}

// Server holds both HTTP and gRPC servers with their configurations
type Server struct {
	HTTP         *http.Server
	GRPC         GRPCServer
	Metrics      *metrics.Metrics
	Tracer       *tracer.Tracer
	Config       *config.Config
	ServerConfig ServerConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	HTTPPort    string
	GRPCPort    string
	MetricsPort string
}

// NewServer creates a new server instance with both HTTP and gRPC servers
func NewServer(cfg *config.Config, httpHandler http.Handler, grpcServer GRPCServer, metricsInstance *metrics.Metrics) (*Server, error) {

	// Initialize tracer
	ctx := context.Background()
	tracerInstance, err := tracer.NewTracer(ctx, cfg.MetricsProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Create HTTP server with metrics middleware
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.HTTPPort),
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		HTTP:    httpServer,
		GRPC:    grpcServer,
		Metrics: metricsInstance,
		Tracer:  tracerInstance,
		Config:  cfg,
		ServerConfig: ServerConfig{
			HTTPPort:    cfg.Server.HTTPPort,
			GRPCPort:    cfg.Server.GRPCPort,
			MetricsPort: cfg.MetricsProvider.PrometheusPort,
		},
	}, nil
}

// Start starts both HTTP and gRPC servers
func (s *Server) Start(ctx context.Context) error {
	// Start system metrics collection
	go func() {
		if err := s.Metrics.StartPrometheusServer(ctx, s.ServerConfig.MetricsPort); err != nil {
			slog.Error("Prometheus metrics server error", "error", err)
		}
	}()

	// Start gRPC server
	go func() {
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.ServerConfig.GRPCPort))
		if err != nil {
			slog.Error("Failed to create gRPC listener", "error", err)
			return
		}

		slog.Info("Starting gRPC server", "port", s.ServerConfig.GRPCPort)
		if err := s.GRPC.Serve(grpcListener); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	// Start HTTP server
	go func() {
		slog.Info("Starting HTTP server", "port", s.ServerConfig.HTTPPort)
		if err := s.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	slog.Info("All servers started successfully")
	return nil
}

// Shutdown gracefully shuts down both servers
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Starting graceful shutdown")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.HTTP.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}

	// Shutdown gRPC server
	s.GRPC.GracefulStop()

	// Shutdown tracer
	if err := s.Tracer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Tracer shutdown error", "error", err)
	}

	slog.Info("All servers shut down successfully")
	return nil
}
