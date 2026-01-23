package router

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	server       *http.Server
	log          *zap.SugaredLogger
	closeTimeout time.Duration
	ctx          context.Context
}

func NewServer(
	ctx context.Context,
	h http.Handler,
	log *zap.SugaredLogger,
	address string,
	closeTimeout time.Duration,
) *Server {
	s := &Server{
		server: &http.Server{
			Addr:        address,
			BaseContext: func(l net.Listener) context.Context { return ctx },
			Handler:     h,
		},
		log:          log,
		closeTimeout: closeTimeout,
		ctx:          ctx,
	}
	return s
}

func (s *Server) Start() error {
	s.log.Infow("Main server is starting...", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Warnw("Main server ListenAndServe()", "err", err)
		return err
	}
	return nil
}

func (s *Server) StartTLS(certFile string, keyFile string) error {
	s.log.Infow("Main server is starting (TLS)...", "addr", s.server.Addr)
	if err := s.server.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Warnw("Main server ListenAndServeTLS()", "err", err)
		return err
	}
	return nil
}

func (s *Server) Close() error {
	ctx, close := context.WithTimeout(s.ctx, s.closeTimeout)
	defer close()
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Warnw("Main server Shutdown()", "err", err)
		return err
	}
	s.log.Info("Main server is closed")
	return nil
}
