package service

import (
	"context"
	"sync"
	"time"

	"github.com/domurdoc/shortener/internal/model"
)

type deleteRes struct {
	count int
	err   error
}

func (s *Service) DeleteShortCodes(ctx context.Context, user *model.User, shortCodes []string) {
	for _, shortCode := range shortCodes {
		select {
		case <-s.doneCh:
			return
		case s.deletedRecords <- model.UserRecord{UserID: user.ID, ShortCode: model.ShortCode(shortCode)}:
		}
	}
}

func (s *Service) serveDeletions() {
	batchCh := s.generator()
	resChs := s.fanOut(batchCh)
	resCh := s.fanIn(resChs...)
	for res := range resCh {
		if res.err != nil {
			s.log.Errorw("failed to process deletions", "err", res.err)
			continue
		}
		s.log.Debugw("deletions saved", "count", res.count)
	}
}

func (s *Service) generator() chan []model.UserRecord {
	batchCh := make(chan []model.UserRecord)

	go func() {
		defer close(batchCh)

		var batch []model.UserRecord

		t := time.NewTicker(s.checkInterval)
		defer t.Stop()

		for {
			select {
			case <-s.doneCh:
				return
			case record := <-s.deletedRecords:
				batch = append(batch, record)
				if len(batch) >= s.maxBatchSize {
					select {
					case <-s.doneCh:
						return
					case batchCh <- batch:
						batch = nil
					}
					t.Reset(s.checkInterval)
				}
			case <-t.C:
				if len(batch) > 0 {
					select {
					case <-s.doneCh:
						return
					case batchCh <- batch:
						batch = nil
					}
				}
			}
		}

	}()

	return batchCh
}

func (s *Service) fanOut(batchCh chan []model.UserRecord) []chan deleteRes {
	resChs := make([]chan deleteRes, s.maxWorkers)
	for i := range s.maxWorkers {
		resChs[i] = s.delete(batchCh)
	}
	return resChs
}

func (s *Service) delete(batchCh chan []model.UserRecord) chan deleteRes {
	resCh := make(chan deleteRes)

	go func() {
		defer close(resCh)

		for batch := range batchCh {
			count, err := s.repo.Delete(context.Background(), batch)
			select {
			case <-s.doneCh:
				return
			case resCh <- deleteRes{count: count, err: err}:
			}
		}
	}()

	return resCh
}

func (s *Service) fanIn(resChs ...chan deleteRes) chan deleteRes {
	finalCh := make(chan deleteRes)

	var wg sync.WaitGroup
	for _, ch := range resChs {
		wg.Add(1)

		go func(ch chan deleteRes) {
			defer wg.Done()

			for res := range ch {
				select {
				case <-s.doneCh:
					return
				case finalCh <- res:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
