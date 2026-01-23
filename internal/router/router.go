package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/logger"
)

func NewRouter(h *handler.Handler, a *auth.Auth, log *zap.SugaredLogger, statsTrustedSubnet string) http.Handler {
	r := chi.NewRouter()

	r.Use(logger.NewRequestLogger(log))
	r.Use(compressor.GZIPMiddleware)

	r.Post("/", auth.AuthRequest(h.Shorten, a))
	r.Get("/ping", auth.AuthRequest(h.Ping, a))
	r.Get("/{shortCode}", auth.AuthRequest(h.Retrieve, a))
	r.Post("/api/shorten", auth.AuthRequest(h.ShortenJSON, a))
	r.Post("/api/shorten/batch", auth.AuthRequest(h.ShortenBatchJSON, a))
	r.Get("/api/user/urls", auth.AuthRequest(h.RetrieveForUser, a))
	r.Delete("/api/user/urls", auth.AuthRequest(h.DeleteShortCodes, a))
	r.Get("/api/internal/stats", httputil.CheckSubnet(h.GetStats, statsTrustedSubnet))

	return r
}
