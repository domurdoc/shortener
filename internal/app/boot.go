package app

import (
	"database/sql"
	"net/http"

	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/router"
	"github.com/domurdoc/shortener/internal/service"
	"go.uber.org/zap"
)

func Boot(logLevel, fileStoragePath, baseURL, databaseDSN string) (http.Handler, *zap.SugaredLogger, *sql.DB) {
	log := logger.New(logLevel)
	db := db.NewPGX(databaseDSN)
	serializer := repository.NewJSONSerializer()
	repo := repository.NewFileRepo(fileStoragePath, serializer)
	service := service.New(repo, baseURL)
	handler := handler.New(service, db)
	router := router.New(handler, log)
	return router, log, db
}
