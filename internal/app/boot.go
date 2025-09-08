package app

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/router"
	"github.com/domurdoc/shortener/internal/service"
	"go.uber.org/zap"
)

func Boot(logLevel, fileStoragePath, baseURL string) (http.Handler, *zap.SugaredLogger) {
	log := logger.New(logLevel)
	serializer := repository.NewJSONSerializer()
	repo := repository.NewFileRepo(fileStoragePath, serializer)
	service := service.New(repo, baseURL)
	handler := handler.New(service)
	router := router.New(handler, log)
	return router, log
}
