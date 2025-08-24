package restclientexample

import (
	"fmt"
	"go-platform/internal/models/dogs"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"resty.dev/v3"
)

const timeout = 10 * time.Second

// DogAPI represents a client for the Dog API
type DogAPI struct {
	rClient *resty.Client
}

// NewDogAPI creates a new Dog API client
func NewDogAPI() *DogAPI {
	return &DogAPI{
		rClient: resty.New().
			SetTimeout(timeout).
			SetBaseURL("https://dog.ceo/api").
			SetHeader("Content-Type", "application/json"),
	}
}

// GetRandomDogImageByBreed gets a random dog image for a specific breed
func (d *DogAPI) GetRandomDogImageByBreed(breed string) (string, error) {
	var response dogs.DogResponse

	res, err := d.rClient.R().
		SetPathParam("breed", breed).
		SetResult(&response).
		Get("/breed/{breed}/images/random")
	if err != nil {
		return "", errors.Wrap(err, "failed to send request")
	}

	if res.IsError() {
		slog.Error("Failed to get random dog image", "breed", breed, "status", res.StatusCode())
		return "", fmt.Errorf("received non-200 response status: %d", res.StatusCode())
	}

	if response.Status != "success" {
		return "", fmt.Errorf("API returned error status: %s", response.Status)
	}

	slog.Info("Successfully retrieved random dog image", "breed", breed)
	return response.Message, nil
}

// DownloadDogImage downloads the actual image from the given URL and returns it as bytes
func (d *DogAPI) DownloadDogImage(imageURL string) ([]byte, error) {
	// Create a new client for downloading images (without base URL)
	imageClient := resty.New().
		SetTimeout(timeout).
		SetHeader("User-Agent", "Go-Platform/1.0")

	res, err := imageClient.R().
		Get(imageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download image")
	}

	if res.IsError() {
		slog.Error("Failed to download image", "url", imageURL, "status", res.StatusCode())
		return nil, fmt.Errorf("received non-200 response status: %d", res.StatusCode())
	}

	imageBytes := res.Bytes()
	slog.Info("Successfully downloaded image", "url", imageURL, "size_bytes", len(imageBytes))

	return imageBytes, nil
}
