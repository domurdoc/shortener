package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/domurdoc/shortener/internal/handler"
)

func New(handler *handler.Handler) http.Handler {
	router := chi.NewRouter()
	setupRoutes(router, handler)
	return router
}

func setupRoutes(router *chi.Mux, handler *handler.Handler) {
	router.Post("/", handler.Shorten)
	router.Get("/ping", handler.Ping)
	router.Get("/{shortCode}", handler.Retrieve)
	router.Post("/api/shorten", handler.ShortenJSON)
	router.Post("/api/shorten/batch", handler.ShortenBatchJSON)
	router.Get("/api/user/urls", handler.RetrieveForUser)
	router.Delete("/api/user/urls", handler.DeleteShortCodes)
}
