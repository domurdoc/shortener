package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

// 2048 - max url length (RFC)
const URLMaxLength = 2048

type Shortener struct {
	repo    repository.RecordRepo
	baseURL string
}

func New(repo repository.RecordRepo, baseURL string) *Shortener {
	return &Shortener{repo: repo, baseURL: baseURL}
}

func (s *Shortener) Shorten(ctx context.Context, user *model.User, originalURL string) (string, error) {
	shortCode, shortURL, err := s.generateShortCodeURL(originalURL)
	if err != nil {
		return "", err
	}
	record := &model.Record{
		OriginalURL: model.OriginalURL(originalURL),
		ShortCode:   model.ShortCode(shortCode),
	}
	err = s.repo.Store(ctx, record, user.ID)
	var urlErr *model.OriginalURLExistsError
	if errors.As(err, &urlErr) {
		shortURL, err := url.JoinPath(s.baseURL, string(urlErr.ShortCode))
		if err != nil {
			return "", err
		}
		return shortURL, ErrURLConflict
	}
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *Shortener) GetByShortCode(ctx context.Context, shortCode string) (string, error) {
	record, err := s.repo.Fetch(ctx, model.ShortCode(shortCode))
	var e *model.ShortCodeNotFoundError
	if errors.As(err, &e) {
		return "", &NotFoundError{shortCode: shortCode}
	}
	if err != nil {
		return "", err
	}
	return string(record.OriginalURL), nil
}

func (s *Shortener) ShortenBatch(ctx context.Context, user *model.User, originalURLS []string) ([]string, error) {
	shortURLS := make([]string, 0, len(originalURLS))
	records := make([]model.Record, 0, len(originalURLS))
	for _, originalURL := range originalURLS {
		shortCode, shortURL, err := s.generateShortCodeURL(originalURL)
		if err != nil {
			return nil, err
		}
		record := model.Record{
			OriginalURL: model.OriginalURL(originalURL),
			ShortCode:   model.ShortCode(shortCode),
		}

		records = append(records, record)
		shortURLS = append(shortURLS, shortURL)
	}
	err := s.repo.StoreBatch(ctx, records, user.ID)
	var batchErr model.BatchError
	if errors.As(err, &batchErr) {
		for _, e := range batchErr {
			shortURL, err := url.JoinPath(s.baseURL, string(e.ShortCode))
			if err != nil {
				return nil, err
			}
			shortURLS[e.BatchPos] = shortURL
		}
		return shortURLS, ErrURLConflict
	}
	if err != nil {
		return nil, err
	}
	return shortURLS, nil
}

func (s *Shortener) GetForUser(ctx context.Context, user *model.User) ([]model.URLRecord, error) {
	records, err := s.repo.FetchForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	urlRecords := make([]model.URLRecord, 0, len(records))
	for _, record := range records {
		shortURL, err := url.JoinPath(s.baseURL, string(record.ShortCode))
		if err != nil {
			return nil, err
		}
		urlRecord := model.URLRecord{
			OriginalURL: record.OriginalURL,
			ShortURL:    model.ShortURL(shortURL),
		}
		urlRecords = append(urlRecords, urlRecord)
	}
	return urlRecords, nil
}

func (s *Shortener) generateShortCodeURL(originalURL string) (string, string, error) {
	if err := validateURL(originalURL); err != nil {
		return "", "", err
	}
	shortCode := generateShortCode()
	shortURL, err := url.JoinPath(s.baseURL, shortCode)
	if err != nil {
		return "", "", err
	}
	return shortCode, shortURL, nil
}

func validateURL(URL string) error {
	if len(URL) > URLMaxLength {
		return &URLError{msg: "url too long", url: URL}
	}
	parsedLongURL, err := url.Parse(URL)
	if err != nil {
		return &URLError{msg: err.Error(), url: URL}
	}
	if parsedLongURL.Host == "" {
		return &URLError{msg: "must be absolute", url: URL}
	}
	if parsedLongURL.String() != URL {
		return &URLError{msg: "must be url-encoded", url: URL}
	}
	return nil
}
