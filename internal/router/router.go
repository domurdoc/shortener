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

func New(handler *handler.Handler, log *zap.SugaredLogger) http.Handler {
	router := chi.NewRouter()
	setupRoutes(router, handler)
	return setupMiddleware(router, log)
}

func setupRoutes(router *chi.Mux, handler *handler.Handler) {
	router.Post("/", handler.Shorten)
	router.Get("/ping", handler.Ping)
	router.Get("/{shortCode}", handler.Retrieve)
	router.Post("/api/shorten", handler.ShortenJSON)
	router.Post("/api/shorten/batch", handler.ShortenBatchJSON)
}

func setupMiddleware(router *chi.Mux, log *zap.SugaredLogger) http.Handler {
	return httputil.AddMiddlewares(
		router,
		logger.NewRequestLogger(log),
		compressor.GZIPMiddleware,
	)
}
