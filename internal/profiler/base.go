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

func (p *ProfilerServer) Start() {
	if err := p.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.log.Fatalw("ProfilerServer.ListenAndServe()", "err", err)
	}
}

func (p *ProfilerServer) Close() error {
	ctx, close := context.WithTimeout(p.ctx, p.closeTimeout)
	defer close()
	return p.server.Shutdown(ctx)
}
