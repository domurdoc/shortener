package repository

import (
	"encoding/json"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"sync"
)

type record struct {
	ID    int   `json:"UUID"`
	Key   Key   `json:"short_url"`
	Value Value `json:"original_url"`
}

func (r record) MarshalJSON() ([]byte, error) {
	type recordAlias record

	aliasValue := struct {
		recordAlias
		ID string `json:"UUID"`
	}{
		recordAlias: recordAlias(r),
		ID:          strconv.Itoa(r.ID),
	}
	return json.Marshal(aliasValue)
}

func (r *record) UnmarshalJSON(data []byte) (err error) {
	type recordAlias record

	aliasValue := &struct {
		*recordAlias
		ID string `json:"UUID"`
	}{
		recordAlias: (*recordAlias)(r),
	}
	if err = json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	r.ID, err = strconv.Atoi(aliasValue.ID)
	return
}

type FileRepo struct {
	path string
	mu   sync.Mutex
}

func NewFileRepo(path string) *FileRepo {
	return &FileRepo{path: path}
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
		storage[key] = record{
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

func (s *FileRepo) load() (map[Key]record, error) {
	storage := make(map[Key]record)
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return storage, nil
	}
	var records []record
	if err := json.Unmarshal(content, &records); err != nil {
		return nil, err
	}
	for _, r := range records {
		storage[r.Key] = r
	}
	return storage, nil
}

func (s *FileRepo) dump(storage map[Key]record) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	records := slices.Collect(maps.Values(storage))
	content, err := json.Marshal(records)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}

func nextSeq(storage map[Key]record) int {
	maxID := 0
	for _, r := range storage {
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	return maxID + 1
}
