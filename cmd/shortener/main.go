package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/router"
)

func main() {
	cfg := config.New()
	if err := config.ParseEnv(cfg); err != nil {
		log.Fatal(err)
	}
	config.ParseArgs(cfg)
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close(nil)
	a.Log.Infow(
		"starting server",
		"addr", a.Config.Server.Address,
		"baseURL", a.Config.Service.BaseURL,
		"logLevel", a.Config.Logger.Level,
		"fileStoragePath", a.Config.Repositories.File.Path,
		"databaseDSN", a.Config.Repositories.DB.DSN,
		"repo", fmt.Sprintf("%T", a.RecordRepo),
		"fileSub", a.AuditFileSub,
		"RemoteSub", a.AuditRemoteSub,
		"ProfileServer", a.Config.Profiler.Address,
	)
	handler := handler.New(a)
	router := router.New(handler)
	router = httputil.AddMiddlewares(
		router,
		logger.NewRequestLogger(a.Log),
		auth.NewAuthMiddleware(a.Auth),
		compressor.GZIPMiddleware,
	)
	log.Fatal(http.ListenAndServe(a.Config.Server.Address, router))
}
