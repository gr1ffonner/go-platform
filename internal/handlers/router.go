package handlers

import (
	"net/http"

	_ "go-platform/api"
	"go-platform/pkg/metrics"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter(h *Handler, httpMetrics *metrics.HTTPMetrics) *mux.Router {
	router := mux.NewRouter()

	// Add metrics middleware first (to capture all requests)
	router.Use(MetricsMiddleware(httpMetrics))

	// Add logging middleware
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
