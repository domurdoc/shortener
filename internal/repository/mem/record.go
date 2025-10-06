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
	ShortCodeRecords   map[model.ShortCode]model.BaseRecord
	ShortCodeUserIDS   map[model.ShortCode]map[model.UserID]model.BaseRecord
	UserIDRecords      map[model.UserID]map[model.ShortCode]model.BaseRecord
	OriginalURLRecords map[model.OriginalURL]model.BaseRecord
	mu                 sync.Mutex
}

func NewMemRecordRepo() *MemRecordRepo {
	return &MemRecordRepo{
		ShortCodeRecords:   make(map[model.ShortCode]model.BaseRecord),
		ShortCodeUserIDS:   make(map[model.ShortCode]map[model.UserID]model.BaseRecord),
		UserIDRecords:      make(map[model.UserID]map[model.ShortCode]model.BaseRecord),
		OriginalURLRecords: make(map[model.OriginalURL]model.BaseRecord),
	}
}

func (r *MemRecordRepo) Store(ctx context.Context, record *model.BaseRecord, userID model.UserID) error {
	err := r.StoreBatch(ctx, []model.BaseRecord{*record}, userID)
	var batchURLExistsErr model.BatchOriginalURLExistsError
	if errors.As(err, &batchURLExistsErr) {
		return batchURLExistsErr[0]
	}
	return err
}

func (r *MemRecordRepo) StoreBatch(ctx context.Context, records []model.BaseRecord, userID model.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range records {
		if _, exists := r.ShortCodeRecords[record.ShortCode]; exists {
			return fmt.Errorf("ShortCode %s already exists", record.ShortCode)
		}
	}
	var batchURLExistsErr model.BatchOriginalURLExistsError
	for pos, record := range records {
		existingRecord, exists := r.OriginalURLRecords[record.OriginalURL]
		if !exists {
			r.OriginalURLRecords[record.OriginalURL] = record
			r.ShortCodeRecords[record.ShortCode] = record
			r.ShortCodeUserIDS[record.ShortCode] = make(map[model.UserID]model.BaseRecord)
		} else if existingRecord.ShortCode != record.ShortCode {
			record.ShortCode = existingRecord.ShortCode
			urlExistsErr := &model.OriginalURLExistsError{
				OriginalURL: record.OriginalURL,
				ShortCode:   existingRecord.ShortCode,
				BatchPos:    pos,
			}
			batchURLExistsErr = append(batchURLExistsErr, urlExistsErr)
		}
		if _, ok := r.UserIDRecords[userID]; !ok {
			r.UserIDRecords[userID] = make(map[model.ShortCode]model.BaseRecord)
		}
		r.UserIDRecords[userID][record.ShortCode] = record
		r.ShortCodeUserIDS[record.ShortCode][userID] = record
	}
	if len(batchURLExistsErr) != 0 {
		return batchURLExistsErr
	}
	return nil
}

func (r *MemRecordRepo) Fetch(ctx context.Context, shortCode model.ShortCode) (*model.BaseRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, exists := r.ShortCodeRecords[shortCode]
	if !exists {
		return nil, &model.ShortCodeNotFoundError{ShortCode: shortCode}
	}
	userIDS := r.ShortCodeUserIDS[shortCode]
	if len(userIDS) == 0 {
		return nil, &model.ShortCodeDeletedError{ShortCode: shortCode}
	}
	return &record, nil
}

func (r *MemRecordRepo) FetchForUser(ctx context.Context, userID model.UserID) ([]model.BaseRecord, error) {
	originalURLRecords, ok := r.UserIDRecords[userID]
	if !ok {
		return nil, nil
	}
	return slices.Collect(maps.Values(originalURLRecords)), nil
}

func (r *MemRecordRepo) Delete(ctx context.Context, records []model.UserRecord) error {
	for _, record := range records {
		userIDS, ok := r.ShortCodeUserIDS[record.ShortCode]
		if !ok {
			continue
		}
		delete(userIDS, record.UserID)
		shortCodeRecords, ok := r.UserIDRecords[record.UserID]
		if !ok {
			continue
		}
		delete(shortCodeRecords, record.ShortCode)
	}
	return nil
}
