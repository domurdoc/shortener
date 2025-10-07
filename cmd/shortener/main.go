package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/router"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	a.Log.Infow(
		"starting server",
		"addr", a.Options.Addr,
		"baseURL", a.Options.BaseURL,
		"logLevel", a.Options.LogLevel,
		"fileStoragePath", a.Options.FileStoragePath,
		"databaseDSN", a.Options.DatabaseDSN,
		"repo", fmt.Sprintf("%T", a.RecordRepo),
	)
	handler := handler.New(a.Service)
	router := router.New(handler)
	router = httputil.AddMiddlewares(
		router,
		logger.NewRequestLogger(a.Log),
		auth.NewAuthMiddleware(a.Auth),
		compressor.GZIPMiddleware,
	)
	log.Fatal(http.ListenAndServe(a.Options.Addr.String(), router))
}
