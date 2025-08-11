package main

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/config"
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
	options := config.ParseArgs()
	h := handler.New(options.BaseURL.String(), service.New(repository.NewMem(), nil))
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", h.Shorten)
	r.Get("/{shortCode}", h.Retrieve)
	return http.ListenAndServe(options.Addr.String(), r)
}
