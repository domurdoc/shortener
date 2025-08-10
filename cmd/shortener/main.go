package main

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	shortenerHandler := handler.New("http://localhost:8080", service.New(repository.NewMem(), nil))
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", shortenerHandler.Shorten)
	mux.HandleFunc("GET /{shortCode}", shortenerHandler.Retrieve)
	return http.ListenAndServe(":8080", mux)
}
