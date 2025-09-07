package router

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/handler"
	"github.com/go-chi/chi/v5"
)

func NewChi(handler *handler.Shortener) chi.Router {
	router := chi.NewRouter()
	router.Post("/", handler.Shorten)
	router.Get("/{shortCode}", handler.GetByShortCode)
	router.Post("/api/shorten", handler.ShortenJSON)
	return router
}

func NewBase(handler *handler.Shortener) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handler.Shorten)
	mux.HandleFunc("GET /{shortCode}", handler.GetByShortCode)
	mux.HandleFunc("POST /api/shorten", handler.ShortenJSON)
	return mux
}
