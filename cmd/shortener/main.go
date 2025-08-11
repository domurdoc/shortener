package main

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	h := handler.New("http://localhost:8080", service.New(repository.NewMem(), nil))
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", h.Shorten)
	r.Get("/{shortCode}", h.Retrieve)
	return http.ListenAndServe(":8080", r)
}
