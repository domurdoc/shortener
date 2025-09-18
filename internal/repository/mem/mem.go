package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/domurdoc/shortener/internal/repository"
)

type MemRepo struct {
	kv map[repository.Key]repository.Value
	vk map[repository.Value]repository.Key
	mu sync.Mutex
}

func New() *MemRepo {
	return &MemRepo{kv: make(map[repository.Key]repository.Value), vk: make(map[repository.Value]repository.Key)}
}

func (s *MemRepo) Store(ctx context.Context, key repository.Key, value repository.Value) error {
	err := s.StoreBatch(ctx, repository.SingleItemBatch(key, value))
	var e repository.BatchError
	if errors.As(err, &e) {
		return (e)[0]
	}
	return err
}

func (s *MemRepo) StoreBatch(ctx context.Context, batchItems []repository.BatchItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range batchItems {
		if _, exists := s.kv[item.Key]; exists {
			return &repository.KeyAlreadyExistsError{Key: item.Key}
		}
	}
	var batchError repository.BatchError
	for pos, item := range batchItems {
		returnedKey, exists := s.vk[item.Value]
		if !exists {
			s.vk[item.Value] = item.Key
			s.kv[item.Key] = item.Value
			continue
		}
		if returnedKey != item.Key {
			valueErr := &repository.ValueAlreadyExistsError{Key: returnedKey, Value: item.Value, Pos: pos}
			batchError = append(batchError, valueErr)
		}
	}
	if len(batchError) != 0 {
		return batchError
	}
	return nil
}

func (s *MemRepo) Fetch(ctx context.Context, key repository.Key) (repository.Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.kv[key]
	if !exists {
		return "", &repository.KeyNotFoundError{Key: key}
	}
	return value, nil
}

func (s *MemRepo) Ping(ctx context.Context) error {
	return nil
}

func (s *MemRepo) Close() error {
	return nil
}
