package subscribers

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"resty.dev/v3"

	"github.com/domurdoc/shortener/internal/audit"
)

type RemoteSubscriber struct {
	id       string
	url      string
	events   chan *audit.Event
	poolSize int
	wg       sync.WaitGroup
	doneCh   chan struct{}
	log      *zap.SugaredLogger
	client   *resty.Client
}

func NewRemote(url string, log *zap.SugaredLogger, poolSize int) *RemoteSubscriber {
	s := &RemoteSubscriber{
		id:       "RemoteSubscriber_" + url,
		url:      url,
		events:   make(chan *audit.Event),
		poolSize: poolSize,
		wg:       sync.WaitGroup{},
		doneCh:   make(chan struct{}),
		log:      log,
		client:   resty.New(),
	}
	for i := range poolSize {
		s.wg.Add(1)
		go s.worker(i)
	}
	return s
}

func (s *RemoteSubscriber) Update(e *audit.Event) error {
	select {
	case <-s.doneCh:
		return fmt.Errorf("closed")
	case s.events <- e:
		return nil
	}
}

func (s *RemoteSubscriber) GetID() string {
	return s.id
}

func (s *RemoteSubscriber) Close() error {
	close(s.doneCh)
	s.wg.Wait()
	return nil
}

func (s *RemoteSubscriber) worker(workerID int) {
	defer s.wg.Done()
	ctx := context.Background()

	for {
		select {
		case <-s.doneCh:
			return
		case e := <-s.events:
			if err := s.write(ctx, e); err != nil {
				s.log.Errorw("failed to write event", "saver_id", s.id, "worker_id", workerID, "error", err.Error())
				continue
			}
		}
	}
}

func (s *RemoteSubscriber) write(_ context.Context, e *audit.Event) error {
	resp, err := s.client.R().SetBody(e).Post(s.url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("bad status code: %d", resp.StatusCode())
	}
	return nil
}
