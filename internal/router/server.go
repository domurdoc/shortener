package router

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/compressor"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/logger"
)

type Server struct {
	server       *http.Server
	log          *zap.SugaredLogger
	closeTimeout time.Duration
	ctx          context.Context
}

func NewServer(
	ctx context.Context,
	a *app.App,
	log *zap.SugaredLogger,
	address string,
	closeTimeout time.Duration,
) *Server {
	h := handler.New(a)
	r := httputil.AddMiddlewares(
		New(h),
		logger.NewRequestLogger(a.Log),
		auth.NewAuthMiddleware(a.Auth),
		compressor.GZIPMiddleware,
	)
	s := &Server{
		server: &http.Server{
			Addr:        address,
			BaseContext: func(l net.Listener) context.Context { return ctx },
			Handler:     r,
		},
		log:          log,
		closeTimeout: closeTimeout,
		ctx:          ctx,
	}
	return s
}

func (s *Server) Start() {
	s.log.Infow("Shortener server is starting...", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalw("ListenAndServe()", "err", err)
	}
}

func (s *Server) StartTLS(certFile string, keyFile string) {
	s.log.Infow("Shortener server is starting (TLS)...", "addr", s.server.Addr)
	if err := s.server.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalw("ListenAndServeTLS()", "err", err)
	}
}

func (s *Server) Close() error {
	ctx, close := context.WithTimeout(s.ctx, s.closeTimeout)
	defer close()
	return s.server.Shutdown(ctx)
}
