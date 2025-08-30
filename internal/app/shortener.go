package app

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/router"
	"github.com/domurdoc/shortener/internal/service"
)

func Run() error {
	options := config.LoadOptions()
	if err := logger.Initialize(options.LogLevel.String()); err != nil {
		return err
	}
	defer logger.Log.Sync()
	repo := repository.NewMem()
	service := service.New(repo)
	handler := handler.New(options.BaseURL.String(), service)
	router := router.NewChi(handler)
	logger.Sugar.Infow(
		"Starting server",
		"addr", options.Addr,
		"baseURL", options.BaseURL,
		"logLevel", options.LogLevel,
	)
	if err := http.ListenAndServe(options.Addr.String(), router); err != nil {
		logger.Sugar.Errorw(
			err.Error(),
			"event", "start server",
		)
		return err
	}
	return nil
}
