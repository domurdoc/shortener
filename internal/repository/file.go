package repository

import (
	"io"
	"maps"
	"os"
	"slices"
	"sync"
)

type FileRepo struct {
	path       string
	mu         sync.Mutex // TODO: use file lock
	serializer Serializer
}

func NewFileRepo(path string, serializer Serializer) *FileRepo {
	repo := FileRepo{path: path, serializer: serializer}
	if _, err := repo.load(); err != nil {
		panic(err)
	}
	return &repo
}

func (s *FileRepo) Store(key Key, value Value) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	storage, err := s.load()
	if err != nil {
		return err
	}
	r, exists := storage[key]
	if !exists {
		storage[key] = Record{
			ID:    nextSeq(storage),
			Key:   key,
			Value: value,
		}
		if err := s.dump(storage); err != nil {
			return err
		}
		return nil
	}
	if r.Value == value {
		return nil
	}
	return &KeyAlreadyExistsError{key: key}
}

func (s *FileRepo) Fetch(key Key) (Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	storage, err := s.load()
	if err != nil {
		return "", err
	}
	r, exists := storage[key]
	if !exists {
		return "", &KeyNotFoundError{key: key}
	}
	return r.Value, nil
}

func (s *FileRepo) load() (map[Key]Record, error) {
	storage := make(map[Key]Record)
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	records, err := s.serializer.Load(content)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		storage[r.Key] = r
	}
	return storage, nil
}

func (s *FileRepo) dump(storage map[Key]Record) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	records := slices.Collect(maps.Values(storage))
	content, err := s.serializer.Dump(records)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}

func nextSeq(storage map[Key]Record) int {
	maxID := 0
	for _, r := range storage {
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	return maxID + 1
}
