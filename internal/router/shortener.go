package router

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(handler *handler.Shortener, log *zap.SugaredLogger) http.Handler {
	router := chi.NewRouter()
	setupRoutes(router, handler)
	return setupMiddleware(router, log)
}

func setupRoutes(router *chi.Mux, handler *handler.Shortener) {
	router.Post("/", handler.Shorten)
	router.Get("/ping", handler.Ping)
	router.Get("/{shortCode}", handler.GetByShortCode)
	router.Post("/api/shorten", handler.ShortenJSON)
}

func setupMiddleware(router *chi.Mux, log *zap.SugaredLogger) http.Handler {
	return httputil.AddMiddlewares(
		router,
		logger.NewRequestLogger(log),
		compressor.GZIPMiddleware,
	)
}
