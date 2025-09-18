package file

import (
	"context"
	"errors"
	"io"
	"maps"
	"os"
	"slices"
	"sync"

	"github.com/domurdoc/shortener/internal/repository"
)

type FileRepo struct {
	path       string
	mu         sync.Mutex // TODO: use file lock
	serializer serializer
}

func New(path string, serializer serializer) *FileRepo {
	repo := FileRepo{path: path, serializer: serializer}
	if err := repo.Ping(context.Background()); err != nil {
		panic(err)
	}
	return &repo
}

type record struct {
	ID    int
	Key   repository.Key
	Value repository.Value
}

type serializer interface {
	Dump([]record) ([]byte, error)
	Load([]byte) ([]record, error)
}

func (s *FileRepo) Store(ctx context.Context, key repository.Key, value repository.Value) error {
	err := s.StoreBatch(ctx, repository.SingleItemBatch(key, value))
	var e repository.BatchError
	if errors.As(err, &e) {
		return (e)[0]
	}
	return err
}

func (s *FileRepo) StoreBatch(ctx context.Context, batchItems []repository.BatchItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	keyStorage, err := s.load()
	if err != nil {
		return err
	}
	valueStorage := toValueStorage(keyStorage)
	for _, item := range batchItems {
		if _, exists := keyStorage[item.Key]; exists {
			return &repository.KeyAlreadyExistsError{Key: item.Key}
		}
	}
	seq := nextSeq(keyStorage)
	var batchError repository.BatchError
	for pos, item := range batchItems {
		returnedRecord, exists := valueStorage[item.Value]
		if !exists {
			record := record{
				ID:    seq,
				Key:   item.Key,
				Value: item.Value,
			}
			valueStorage[item.Value] = record
			keyStorage[item.Key] = record
			seq++
			continue
		}
		if returnedRecord.Key != item.Key {
			valueErr := &repository.ValueAlreadyExistsError{Key: returnedRecord.Key, Value: item.Value, Pos: pos}
			batchError = append(batchError, valueErr)
		}
	}
	if err := s.dump(keyStorage); err != nil {
		return err
	}
	if len(batchError) != 0 {
		return batchError
	}
	return nil
}

func (s *FileRepo) Fetch(ctx context.Context, key repository.Key) (repository.Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	storage, err := s.load()
	if err != nil {
		return "", err
	}
	r, exists := storage[key]
	if !exists {
		return "", &repository.KeyNotFoundError{Key: key}
	}
	return r.Value, nil
}

func (s *FileRepo) Ping(ctx context.Context) error {
	_, err := s.load()
	return err
}

func (s *FileRepo) Close() error {
	return nil
}

func (s *FileRepo) load() (map[repository.Key]record, error) {
	storage := make(map[repository.Key]record)
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

func (s *FileRepo) dump(storage map[repository.Key]record) error {
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

func nextSeq(storage map[repository.Key]record) int {
	maxID := 0
	for _, r := range storage {
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	return maxID + 1
}

func toValueStorage(keyStorage map[repository.Key]record) map[repository.Value]record {
	valueStorage := make(map[repository.Value]record, len(keyStorage))
	for _, r := range keyStorage {
		valueStorage[r.Value] = r
	}
	return valueStorage
}
