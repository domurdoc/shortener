package mem

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/domurdoc/shortener/internal/model"
)

type MemRecordRepo struct {
	ShortCodeRecords   map[model.ShortCode]model.Record
	UserIDRecords      map[model.UserID]map[model.ShortCode]model.Record
	OriginalURLRecords map[model.OriginalURL]model.Record
	mu                 sync.Mutex
}

func NewMemRecordRepo() *MemRecordRepo {
	return &MemRecordRepo{
		ShortCodeRecords:   make(map[model.ShortCode]model.Record),
		UserIDRecords:      make(map[model.UserID]map[model.ShortCode]model.Record),
		OriginalURLRecords: make(map[model.OriginalURL]model.Record),
	}
}

func (s *MemRecordRepo) Store(ctx context.Context, record *model.Record, userID model.UserID) error {
	err := s.StoreBatch(ctx, []model.Record{*record}, userID)
	var batchErr model.BatchError
	if errors.As(err, &batchErr) {
		return batchErr[0]
	}
	return err
}

func (s *MemRecordRepo) StoreBatch(ctx context.Context, records []model.Record, userID model.UserID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, record := range records {
		if _, exists := s.ShortCodeRecords[record.ShortCode]; exists {
			return fmt.Errorf("ShortCode %s already exists", record.ShortCode)
		}
	}
	var batchError model.BatchError
	for pos, record := range records {
		existingRecord, exists := s.OriginalURLRecords[record.OriginalURL]
		if !exists {
			s.OriginalURLRecords[record.OriginalURL] = record
			s.ShortCodeRecords[record.ShortCode] = record
		} else if existingRecord.ShortCode != record.ShortCode {
			record.ShortCode = existingRecord.ShortCode
			urlErr := &model.OriginalURLExistsError{
				OriginalURL: record.OriginalURL,
				ShortCode:   existingRecord.ShortCode,
				BatchPos:    pos,
			}
			batchError = append(batchError, urlErr)
		}
		if _, ok := s.UserIDRecords[userID]; !ok {
			s.UserIDRecords[userID] = make(map[model.ShortCode]model.Record)
		}
		s.UserIDRecords[userID][record.ShortCode] = record
	}
	if len(batchError) != 0 {
		return batchError
	}
	return nil
}

func (s *MemRecordRepo) Fetch(ctx context.Context, shortCode model.ShortCode) (*model.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, exists := s.ShortCodeRecords[shortCode]
	if !exists {
		return nil, &model.ShortCodeNotFoundError{ShortCode: shortCode}
	}
	return &record, nil
}

func (s *MemRecordRepo) FetchForUser(ctx context.Context, userID model.UserID) ([]model.Record, error) {
	originalURLRecords, ok := s.UserIDRecords[userID]
	if !ok {
		return nil, nil
	}
	return slices.Collect(maps.Values(originalURLRecords)), nil
}
