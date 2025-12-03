package profiler

import (
	"context"
	"log"
	"net/http"
	"sync"
)

type Profiler struct {
	Address string
	server  *http.Server
	wg      *sync.WaitGroup
}

func New(address string) *Profiler {
	p := &Profiler{
		Address: address,
		wg:      &sync.WaitGroup{},
	}
	p.server = &http.Server{
		Addr: p.Address,
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		if err := p.server.ListenAndServe(); err != nil {
			log.Fatalf("ProfileServer.ListenAndServe(): %v", err)
		}
	}()
	return p
}

func (p *Profiler) Close() error {
	err := p.server.Shutdown(context.Background())
	p.wg.Wait()
	return err
}
