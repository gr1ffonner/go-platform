package grpc

import (
	"context"
	"log/slog"

	proto "go-platform/api/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ValidationInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch val := req.(type) {
	case *proto.HealthCheckRequest:
		slog.Info("Middleware for HealthCheckRequest started", "request", val)
		return handler(ctx, req)
	}
	return handler(ctx, req)

}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Log start of gRPC call
	slog.InfoContext(
		ctx,
		"gRPC call started",
		"method", info.FullMethod,
	)

	// Execute the handler
	resp, err := handler(ctx, req)

	// Log end of gRPC call
	if err == nil {
		slog.InfoContext(
			ctx,
			"gRPC call completed successfully",
			"method", info.FullMethod,
		)
	} else {
		grpcErr, ok := status.FromError(err)
		if ok {
			slog.ErrorContext(
				ctx,
				"gRPC call failed",
				"method", info.FullMethod,
				"error_msg", grpcErr.Message(),
				"status_code", grpcErr.Code().String(),
			)
		} else {
			slog.ErrorContext(
				ctx,
				"gRPC call failed",
				"method", info.FullMethod,
				"error_msg", err.Error(),
			)
		}
	}

	return resp, err
}
