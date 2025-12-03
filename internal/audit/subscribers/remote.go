package subscribers

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"resty.dev/v3"

	"github.com/domurdoc/shortener/internal/audit"
)

// RemoteSubscriber is a concrete implementation of audit.Subscriber that sends audit events to a remote HTTP endpoint.
// It uses a pool of worker goroutines to concurrently transmit events in the background.
// Events are sent as JSON over HTTP POST requests, and failed attempts are logged but not retried.
// The subscriber ensures graceful shutdown by waiting for active workers to finish.
type RemoteSubscriber struct {
	id       string             // id is the unique identifier for this subscriber.
	url      string             // url is the remote HTTP endpoint to which audit events are sent.
	events   chan *audit.Event  // events is the channel that receives incoming audit events for transmission.
	poolSize int                // poolSize specifies the number of concurrent worker goroutines handling event delivery.
	wg       sync.WaitGroup     // wg tracks the active worker goroutines to ensure all finish during shutdown.
	doneCh   chan struct{}      // doneCh signals the subscriber to stop processing events.
	log      *zap.SugaredLogger // log is the logger used for internal logging, including transmission errors.
	client   *resty.Client      // client is the HTTP client used for sending requests to the remote endpoint.
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
