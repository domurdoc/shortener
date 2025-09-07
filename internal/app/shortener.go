package app

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputils"
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
	repo := repository.NewFileRepo(options.FileStoragePath.String())
	service := service.New(repo, options.BaseURL.String())
	handler := handler.New(service)
	router := router.NewChi(handler)
	middlewares := []httputils.Middleware{
		compressor.GZIPMiddleware,
		logger.RequestLogger,
	}
	logger.Sugar.Infow(
		"Starting server",
		"addr", options.Addr,
		"baseURL", options.BaseURL,
		"logLevel", options.LogLevel,
	)
	if err := http.ListenAndServe(options.Addr.String(), httputils.AddMiddlewares(router, middlewares...)); err != nil {
		logger.Sugar.Errorw(
			err.Error(),
			"event", "start server",
		)
		return err
	}
	return nil
}
