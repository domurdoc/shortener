package app

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/repository"
	dbRepo "github.com/domurdoc/shortener/internal/repository/db"
	fileRepo "github.com/domurdoc/shortener/internal/repository/file"
	memRepo "github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/router"
	"github.com/domurdoc/shortener/internal/service"
	"go.uber.org/zap"
)

func Boot(
	logLevel,
	baseURL,
	databaseDSN,
	fileStoragePath string,
) (
	repository.Repo,
	*zap.SugaredLogger,
	http.Handler,
) {
	log := logger.New(logLevel)
	repo := bootRepo(databaseDSN, fileStoragePath)
	service := service.New(repo, baseURL)
	handler := handler.New(service)
	router := router.New(handler, log)
	return repo, log, router
}

func bootRepo(databaseDSN, fileStoragePath string) repository.Repo {
	if databaseDSN != "" {
		pgDB := db.NewPG(databaseDSN)
		db.MigratePG(pgDB)
		return dbRepo.New(pgDB, dbRepo.NewPGArger)
	}
	if fileStoragePath != "" {
		jsonSerializer := fileRepo.NewJSONSerializer()
		return fileRepo.New(fileStoragePath, jsonSerializer)
	}
	return memRepo.New()
}
