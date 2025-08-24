package handlers

import (
	"log/slog"
	"net/http"

	"go-platform/internal/models/dogs"
	httputils "go-platform/pkg/utils/http-utils"

	"github.com/gorilla/mux"
)

// GetRandomDogImageByBreed godoc
//
//	@Summary		Get random dog image by breed
//	@Description	Retrieves a random dog image for the specified breed, downloads it, and uploads to S3
//	@Tags			Dogs
//	@Param			breed	path	string	true	"Dog breed"
//	@Produce		json
//	@Success		200	{object}	dogs.DogImageResponse	"S3 URL of the uploaded image"
//	@Failure		400	{object}	httputils.ErrorResponse
//	@Failure		404	{object}	httputils.ErrorResponse
//	@Failure		500	{object}	httputils.ErrorResponse
//	@Router			/api/v1/dogs/{breed}/image [get]
func (h *Handler) GetRandomDogImageByBreed(w http.ResponseWriter, r *http.Request) {

	// Extract breed from URL path
	vars := mux.Vars(r)
	breed := vars["breed"]
	if breed == "" {
		slog.Error("Empty breed parameter")
		httputils.WriteResponse(w, http.StatusBadRequest, "Breed parameter is required", nil, nil)
		return
	}
	slog.Info("Processing dog image request", "breed", breed)

	// Call service layer
	imageURL, err := h.dogsService.GetRandomDogImage(r.Context(), breed)
	if err != nil {
		slog.Error("Service failed", "breed", breed, "error", err)
		httputils.WriteResponse(w, http.StatusInternalServerError, "Failed to get dog image", err, nil)
		return
	}
	slog.Info("Service completed", "breed", breed, "image_url", imageURL)

	// Return success response
	response := dogs.DogImageResponse{
		ImageURL: imageURL,
		Breed:    breed,
	}
	httputils.WriteResponse(w, http.StatusOK, "Dog image retrieved successfully", nil, response)
}
