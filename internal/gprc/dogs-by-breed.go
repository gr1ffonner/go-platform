package grpc

import (
	"context"
	proto "go-platform/api/protobuf"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) GetRandomDogImage(ctx context.Context, req *proto.GetRandomDogImageRequest) (*proto.GetRandomDogImageResponse, error) {
	breed := req.GetBreed()

	imageURL, err := s.dogsService.GetRandomDogImage(ctx, breed)
	if err != nil {
		slog.Error("Service failed", "breed", breed, "error", err)
		return nil, status.Errorf(codes.Internal, "Failed to get dog image")
	}
	slog.Info("Service completed", "breed", breed, "image_url", imageURL)

	return &proto.GetRandomDogImageResponse{
		ImageUrl:  imageURL,
		Breed:     breed,
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
