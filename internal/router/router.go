package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/service"
)

func Run(service *service.Shortener, auth *auth.Auth, db *sql.DB, log *zap.SugaredLogger, address string) error {
	handler := handler.New(service, auth, db)
	router := New(handler, log)
	return http.ListenAndServe(address, router)
}

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
	router.Get("/api/user/urls", handler.RetrieveForUser)
	router.Delete("/api/user/urls", handler.DeleteShortCodes)
}

func setupMiddleware(router *chi.Mux, log *zap.SugaredLogger) http.Handler {
	return httputil.AddMiddlewares(
		router,
		logger.NewRequestLogger(log),
		compressor.GZIPMiddleware,
	)
}
