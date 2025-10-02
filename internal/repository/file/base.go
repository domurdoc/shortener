package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"sync"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository/file/serializer"
	"github.com/domurdoc/shortener/internal/repository/mem"
)

type FileRepo struct {
	filepath   string
	serializer serializer.Serializer
	mu         sync.Mutex // TODO: use file lock
}

func New(filepath string, serializer serializer.Serializer) (*FileRepo, error) {
	repo := FileRepo{filepath: filepath, serializer: serializer}
	if _, err := repo.loadMemRepo(context.TODO()); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *FileRepo) Store(ctx context.Context, record *model.Record, userID model.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	memRepo, err := r.loadMemRepo(ctx)
	if err != nil {
		return err
	}
	err = memRepo.Store(ctx, record, userID)
	var urlErr *model.OriginalURLExistsError
	if err != nil && !errors.As(err, &urlErr) {
		return err
	}
	dumpErr := r.dumpMemRepo(memRepo)
	if dumpErr != nil {
		return dumpErr
	}
	return err
}

func (r *FileRepo) StoreBatch(ctx context.Context, records []model.Record, userID model.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	memRepo, err := r.loadMemRepo(ctx)
	if err != nil {
		return err
	}
	err = memRepo.StoreBatch(ctx, records, userID)
	var batchErr model.BatchError
	if err != nil && !errors.As(err, &batchErr) {
		return err
	}
	dumpErr := r.dumpMemRepo(memRepo)
	if dumpErr != nil {
		return dumpErr
	}
	return err
}

func (r *FileRepo) Fetch(ctx context.Context, shortCode model.ShortCode) (*model.Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	memRepo, err := r.loadMemRepo(ctx)
	if err != nil {
		return nil, err
	}
	return memRepo.Fetch(ctx, shortCode)
}

func (r *FileRepo) FetchForUser(ctx context.Context, userID model.UserID) ([]model.Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	memRepo, err := r.loadMemRepo(ctx)
	if err != nil {
		return nil, err
	}
	return memRepo.FetchForUser(ctx, userID)
}

func (r *FileRepo) loadMemRepo(_ context.Context) (*mem.MemRecordRepo, error) {
	memRepo := mem.NewMemRecordRepo()

	file, err := os.OpenFile(r.filepath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	snapshot, err := r.serializer.Load(content)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return memRepo, nil
	}

	shortCodeRecords := make(map[model.ShortCode]model.Record)
	originalURLRecords := make(map[model.OriginalURL]model.Record)
	for _, record := range snapshot.Records {
		shortCodeRecords[record.ShortCode] = record
		originalURLRecords[record.OriginalURL] = record
	}
	userIDRecords := make(map[model.UserID]map[model.ShortCode]model.Record)
	for _, ownership := range snapshot.Ownership {
		record, ok := shortCodeRecords[ownership.ShortCode]
		if !ok {
			return nil, fmt.Errorf("no matching ShortCode")
		}
		if _, ok = userIDRecords[ownership.UserID]; !ok {
			userIDRecords[ownership.UserID] = make(map[model.ShortCode]model.Record)
		}
		userIDRecords[ownership.UserID][record.ShortCode] = record
	}
	memRepo.ShortCodeRecords = shortCodeRecords
	memRepo.UserIDRecords = userIDRecords
	memRepo.OriginalURLRecords = originalURLRecords

	return memRepo, nil
}

func (r *FileRepo) dumpMemRepo(memRepo *mem.MemRecordRepo) error {
	file, err := os.OpenFile(r.filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	records := slices.Collect(maps.Values(memRepo.ShortCodeRecords))
	ownership := make([]serializer.Ownership, 0, len(records))
	for userID, shortCodeRecods := range memRepo.UserIDRecords {
		for shortCode := range shortCodeRecods {
			o := serializer.Ownership{
				UserID:    userID,
				ShortCode: shortCode,
			}
			ownership = append(ownership, o)
		}
	}

	snapshot := &serializer.Snapshot{
		Records:   records,
		Ownership: ownership,
	}

	content, err := r.serializer.Dump(snapshot)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}
