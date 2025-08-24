package grpc

import (
	"net"

	proto "go-platform/api/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	grpcServer *grpc.Server
	proto.UnimplementedHealthServer
}

func NewServer() *server {
	s := &server{
		grpcServer: grpc.NewServer(grpc.ChainUnaryInterceptor(
			LogInterceptor,
			ValidationInterceptor,
		)),
	}

	proto.RegisterHealthServer(s.grpcServer, s)

	reflection.Register(s.grpcServer)

	return s
}

func (s *server) Serve(listener net.Listener) error {
	return s.grpcServer.Serve(listener)
}

func (s *server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
