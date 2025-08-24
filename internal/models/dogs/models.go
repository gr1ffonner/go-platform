package dogs

type DogResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DogImageResponse struct {
	ImageURL string `json:"image_url"`
	Breed    string `json:"breed"`
}
