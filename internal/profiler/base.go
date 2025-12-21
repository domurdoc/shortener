package profiler

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Profiler struct {
	server       *http.Server
	log          *zap.SugaredLogger
	closeTimeout time.Duration
	ctx          context.Context
}

func New(ctx context.Context, address string, log *zap.SugaredLogger, closeTimeout time.Duration) *Profiler {
	return &Profiler{
		server:       &http.Server{Addr: address, BaseContext: func(l net.Listener) context.Context { return ctx }},
		log:          log,
		closeTimeout: closeTimeout,
		ctx:          ctx,
	}
}

func (p *Profiler) Start() {
	if err := p.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.log.Fatalw("Profiler.ListenAndServe()", "err", err)
	}
}

func (p *Profiler) Close() error {
	ctx, close := context.WithTimeout(p.ctx, p.closeTimeout)
	defer close()
	return p.server.Shutdown(ctx)
}
