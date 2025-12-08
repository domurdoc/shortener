package profiler

import (
	"context"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type Profiler struct {
	server *http.Server
	wg     *sync.WaitGroup
	log    *zap.SugaredLogger
}

func New(address string, log *zap.SugaredLogger) *Profiler {
	return &Profiler{
		wg:     &sync.WaitGroup{},
		server: &http.Server{Addr: address},
		log:    log,
	}
}

func (p *Profiler) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := p.server.ListenAndServe(); err != nil {
			p.log.Fatalw("Profiler.ListenAndServe()", "err", err)
		}
	}()
}

func (p *Profiler) Close() error {
	err := p.server.Shutdown(context.Background())
	p.wg.Wait()
	return err
}
