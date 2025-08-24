package handlers

import (
	"net/http"

	_ "go-platform/api"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.Use(LoggingMiddleware)

	// Health
	{
		router.HandleFunc("/live", h.Health).Methods(http.MethodGet)
	}

	// Swagger
	{
		// Redirect /swagger to /swagger/index.html
		router.Handle("/documentation", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently)).Methods(http.MethodGet)

		// Serve Swagger UI
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	}

	// Dogs
	{
		router.HandleFunc("/api/v1/dogs/{breed}/image", h.GetRandomDogImageByBreed).Methods(http.MethodGet)
	}

	return router
}
