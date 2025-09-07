package repository

import (
	"sync"
)

type MemRepo struct {
	storage map[Key]Value
	mu      sync.Mutex
}

func NewMemRepo() *MemRepo {
	return &MemRepo{storage: make(map[Key]Value)}
}

func (s *MemRepo) Store(key Key, value Value) error {
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
	return &KeyAlreadyExistsError{key: key}
}

func (s *MemRepo) Fetch(key Key) (Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.storage[key]
	if !exists {
		return "", &KeyNotFoundError{key: key}
	}
	return value, nil
}
