package mem

import (
	"context"
	"sync"

	"github.com/domurdoc/shortener/internal/repository"
)

type MemRepo struct {
	storage map[repository.Key]repository.Value
	mu      sync.Mutex
}

func New() *MemRepo {
	return &MemRepo{storage: make(map[repository.Key]repository.Value)}
}

func (s *MemRepo) Store(ctx context.Context, key repository.Key, value repository.Value) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	value0, exists := s.storage[key]
	if !exists {
		s.storage[key] = value
		return nil
	}
	if value0 == value {
		return nil
	}
	return &repository.KeyAlreadyExistsError{Key: key}
}

func (s *MemRepo) Fetch(ctx context.Context, key repository.Key) (repository.Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.storage[key]
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
