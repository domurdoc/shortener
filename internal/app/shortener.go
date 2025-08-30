package app

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/router"
	"github.com/domurdoc/shortener/internal/service"
)

func Run() error {
	options := config.LoadOptions()
	repo := repository.NewMem()
	service := service.New(repo)
	handler := handler.New(options.BaseURL.String(), service)
	router := router.NewChi(handler)
	return http.ListenAndServe(options.Addr.String(), router)
}
