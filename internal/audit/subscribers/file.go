package subscribers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/audit"
	"github.com/domurdoc/shortener/internal/utils"
)

type FileSubscriber struct {
	id            string
	filePath      string
	events        chan *audit.Event
	eventBatch    chan []*audit.Event
	poolSize      int
	maxBatchSize  int
	batchInterval time.Duration
	wg            sync.WaitGroup
	doneCh        chan struct{}
	log           *zap.SugaredLogger
	mu            sync.Mutex
}

func NewFile(filePath string, poolSize int, maxBatchSize int, batchInterval time.Duration, log *zap.SugaredLogger) *FileSubscriber {
	s := &FileSubscriber{
		id:            "FileSub_" + utils.GenerateRandomString(utils.ALPHA, 4),
		filePath:      filePath,
		events:        make(chan *audit.Event),
		eventBatch:    make(chan []*audit.Event),
		poolSize:      poolSize,
		maxBatchSize:  maxBatchSize,
		batchInterval: batchInterval,
		wg:            sync.WaitGroup{},
		doneCh:        make(chan struct{}),
		log:           log,
	}
	s.wg.Add(1)
	go s.batcher()
	for i := range s.poolSize {
		s.wg.Add(1)
		go s.worker(i)
	}
	return s
}

func (s *FileSubscriber) Update(e *audit.Event) error {
	select {
	case <-s.doneCh:
		return fmt.Errorf("closed")
	case s.events <- e:
		return nil
	}
}

func (s *FileSubscriber) GetID() string {
	return s.id
}

func (s *FileSubscriber) Close() error {
	close(s.doneCh)
	s.wg.Wait()
	return nil
}

func (s *FileSubscriber) batcher() {
	defer s.wg.Done()
	var batch []*audit.Event

	t := time.NewTicker(s.batchInterval)
	defer t.Stop()

	for {
		select {
		case <-s.doneCh:
			return
		case e := <-s.events:
			s.log.Debugw("received an event to save", "e", e)
			batch = append(batch, e)
			if len(batch) >= s.maxBatchSize {
				select {
				case <-s.doneCh:
					return
				case s.eventBatch <- batch:
					batch = nil
				}
				t.Reset(s.batchInterval)
			}
		case <-t.C:
			if len(batch) > 0 {
				select {
				case <-s.doneCh:
					return
				case s.eventBatch <- batch:
					batch = nil
				}
			}
		}
	}
}

func (s *FileSubscriber) worker(workerID int) {
	defer s.wg.Done()
	ctx := context.Background()

	for {
		select {
		case <-s.doneCh:
			return
		case eventBatch := <-s.eventBatch:
			if err := s.write(ctx, eventBatch); err != nil {
				s.log.Errorw("failed to write events", "saver_id", s.id, "worker_id", workerID, "error", err.Error())
				continue
			}
		}
	}
}

func (s *FileSubscriber) write(_ context.Context, events []*audit.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := serializeEvents(events)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func serializeEvents(events []*audit.Event) ([]byte, error) {
	jsonEvents := make([][]byte, len(events))

	for i, e := range events {
		jsonEvent, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}
		jsonEvents[i] = jsonEvent
	}
	result := bytes.Join(jsonEvents, []byte("\n"))
	result = append(result, '\n')
	return result, nil
}
