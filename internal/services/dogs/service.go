package dogs

import (
	"context"
	"fmt"
	models "go-platform/internal/models/dogs"
	"log/slog"

	"github.com/google/uuid"
)

type DogAPIClient interface {
	GetRandomDogImageByBreed(breed string) (string, error)
	DownloadDogImage(imageURL string) ([]byte, error)
}
type ClientS3 interface {
	PutObject(ctx context.Context, key string, data []byte) (err error)
	GenerateURL(key string) string
}

type Repository interface {
	// return string due to clickhouse dont have auto increment and
	// we should use uuid for simple row
	InsertDog(ctx context.Context, dog *models.Dog) (string, error)
}

type DogsService struct {
	dogAPI     DogAPIClient
	clientS3   ClientS3
	repository Repository
}

func NewDogsService(dogAPI DogAPIClient, clientS3 ClientS3, repository Repository) *DogsService {
	return &DogsService{
		dogAPI:     dogAPI,
		clientS3:   clientS3,
		repository: repository,
	}
}

// GetRandomDogImage gets a random dog image url for a breed and returns the image url
func (s *DogsService) GetRandomDogImage(ctx context.Context, breed string) (string, error) {
	slog.Info("Starting dog image retrieval", "breed", breed)

	// First get the image URL
	imageURL, err := s.dogAPI.GetRandomDogImageByBreed(breed)
	if err != nil {
		slog.Error("Failed to get image URL", "breed", breed, "error", err)
		return "", fmt.Errorf("failed to get image URL: %w", err)
	}
	slog.Info("Got image URL", "breed", breed, "url", imageURL)

	// Then download the actual image
	imageBytes, err := s.dogAPI.DownloadDogImage(imageURL)
	if err != nil {
		slog.Error("Failed to download image", "breed", breed, "url", imageURL, "error", err)
		return imageURL, fmt.Errorf("failed to download image: %w", err)
	}
	slog.Info("Downloaded image", "breed", breed, "size", len(imageBytes))

	// Then upload the image to S3
	imageKey := fmt.Sprintf("dogs/%s/%s", breed, uuid.New().String())
	err = s.clientS3.PutObject(ctx, imageKey, imageBytes)
	if err != nil {
		slog.Error("Failed to upload to S3", "breed", breed, "key", imageKey, "error", err)
		return imageURL, fmt.Errorf("failed to upload to S3: %w", err)
	}
	slog.Info("Uploaded to S3", "breed", breed, "key", imageKey)

	// Generate and return the S3 URL
	s3URL := s.clientS3.GenerateURL(imageKey)
	slog.Info("Dog image retrieval completed", "breed", breed, "s3_url", s3URL)
	return s3URL, nil
}
