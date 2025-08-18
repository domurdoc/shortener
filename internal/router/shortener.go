package router

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewChi(handler *handler.Shortener) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/", handler.Shorten)
	router.Get("/{shortCode}", handler.GetByShortCode)
	return router
}

func NewBase(handler *handler.Shortener) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handler.Shorten)
	mux.HandleFunc("GET /{shortCode}", handler.GetByShortCode)
	return mux
}
