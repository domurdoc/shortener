package profiler

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ProfilerServer struct {
	server       *http.Server
	log          *zap.SugaredLogger
	closeTimeout time.Duration
	ctx          context.Context
}

func NewServer(ctx context.Context, address string, log *zap.SugaredLogger, closeTimeout time.Duration) *ProfilerServer {
	return &ProfilerServer{
		server:       &http.Server{Addr: address, BaseContext: func(l net.Listener) context.Context { return ctx }},
		log:          log,
		closeTimeout: closeTimeout,
		ctx:          ctx,
	}
}

func (p *ProfilerServer) Start() error {
	p.log.Infow("Profiler server is starting...", "addr", p.server.Addr)
	if err := p.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.log.Warnw("Profiler server ListenAndServe()", "err", err)
		return err
	}
	return nil
}

func (p *ProfilerServer) Close() error {
	ctx, close := context.WithTimeout(p.ctx, p.closeTimeout)
	defer close()
	if err := p.server.Shutdown(ctx); err != nil {
		p.log.Warnw("Profiler server Shutdown()", "err", err)
		return err
	}
	p.log.Info("Profiler server is stopped")
	return nil
}
