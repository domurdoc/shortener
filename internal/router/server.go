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
	enableTLS    bool
	certFile     string
	keyFile      string
}

func NewServer(
	ctx context.Context,
	address string,
	enableHTTPS bool,
	certFile string,
	keyFile string,
	closeTimeout time.Duration,
	a *app.App,
	log *zap.SugaredLogger,
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
		enableTLS:    enableHTTPS,
		certFile:     certFile,
		keyFile:      keyFile,
	}
	return s
}

func (s *Server) Start() {
	var serve func() error
	if s.enableTLS {
		serve = func() error { return s.server.ListenAndServeTLS(s.certFile, s.keyFile) }
	} else {
		serve = s.server.ListenAndServe
	}
	if err := serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalw("serve()", "err", err)
	}
}

func (s *Server) Close() error {
	ctx, close := context.WithTimeout(s.ctx, s.closeTimeout)
	defer close()
	return s.server.Shutdown(ctx)
}
