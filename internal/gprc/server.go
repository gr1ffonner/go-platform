package grpc

import (
	"context"
	"net"

	proto "go-platform/api/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DogsService interface {
	GetRandomDogImage(ctx context.Context, breed string) (string, error)
}

type server struct {
	dogsService DogsService
	grpcServer  *grpc.Server
	proto.UnimplementedHealthServer
	proto.UnimplementedDogServiceServer
}

func NewServer(dogsService DogsService) *server {
	s := &server{
		dogsService: dogsService,
		grpcServer: grpc.NewServer(grpc.ChainUnaryInterceptor(
			LogInterceptor,
			ValidationInterceptor,
		)),
	}

	proto.RegisterHealthServer(s.grpcServer, s)
	proto.RegisterDogServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	return s
}

func (s *server) Serve(listener net.Listener) error {
	return s.grpcServer.Serve(listener)
}

func (s *server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
