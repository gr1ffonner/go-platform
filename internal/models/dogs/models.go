package dogs

import "time"

type DogResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DogImageResponse struct {
	ImageURL string `json:"image_url"`
	Breed    string `json:"breed"`
}

type Dog struct {
	ID        string    `json:"id"`
	ImageURL  string    `json:"image_url"`
	Breed     string    `json:"breed"`
	CreatedAt time.Time `json:"created_at"`
}
