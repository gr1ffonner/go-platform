package grpc

import (
	"context"
	proto "go-platform/api/protobuf"
	"log/slog"
)

func (s *server) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	slog.Info("GRPC handler for HealthCheckRequest started", "request", req)
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func (s *server) Watch(req *proto.HealthCheckRequest, stream proto.Health_WatchServer) error {
	return nil
}
