package handlers

import "context"

type DogsService interface {
	GetRandomDogImage(ctx context.Context, breed string) (string, error)
}

type Handler struct {
	dogsService DogsService
}

func NewHandler(dogsService DogsService) *Handler {
	return &Handler{
		dogsService: dogsService,
	}
}
