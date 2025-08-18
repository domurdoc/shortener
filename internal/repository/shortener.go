package repository

import (
	"fmt"
	"sync"
)

type Key string
type Value string

type Shortener interface {
	Store(Key, Value) error
	Fetch(Key) (Value, error)
}

type MemShortener struct {
	storage map[Key]Value
	mu      sync.Mutex
}

func NewMem() *MemShortener {
	return &MemShortener{storage: make(map[Key]Value)}
}

func (s *MemShortener) Store(key Key, value Value) error {
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
	return fmt.Errorf("key %q already exists", key)
}

func (s *MemShortener) Fetch(key Key) (Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.storage[key]
	if !exists {
		return "", fmt.Errorf("key %q not found", key)
	}
	return value, nil
}
