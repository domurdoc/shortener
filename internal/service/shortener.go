package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/domurdoc/shortener/internal/repository"
)

// 2048 - max url length (RFC)
const URLMaxLength = 2048

type Shortener struct {
	repo    repository.Repo
	baseURL string
}

func New(repo repository.Repo, baseURL string) *Shortener {
	return &Shortener{repo: repo, baseURL: baseURL}
}

func (s *Shortener) Shorten(ctx context.Context, longURL string) (string, error) {
	shortCode, shortURL, err := s.generateShortCodeURL(longURL)
	if err != nil {
		return "", err
	}
	return shortURL, s.repo.Store(ctx, repository.Key(shortCode), repository.Value(longURL))
}

func (s *Shortener) GetByShortCode(ctx context.Context, shortCode string) (string, error) {
	url, err := s.repo.Fetch(ctx, repository.Key(shortCode))
	var keyNotFoundError *repository.KeyNotFoundError
	if errors.As(err, &keyNotFoundError) {
		return "", &NotFoundError{shortCode: shortCode}
	}
	if err != nil {
		return "", err
	}
	return string(url), nil
}

func (s *Shortener) ShortenBatch(ctx context.Context, longURLS []string) ([]string, error) {
	shortCodes := make([]string, len(longURLS))
	shortURLS := make([]string, len(longURLS))
	batchItems := make([]repository.BatchItem, len(longURLS))
	for i, longURL := range longURLS {
		shortCode, shortURL, err := s.generateShortCodeURL(longURL)
		if err != nil {
			return nil, err
		}
		shortCodes[i] = shortCode
		shortURLS[i] = shortURL
		batchItems[i] = repository.BatchItem{Key: repository.Key(shortCode), Value: repository.Value(longURL)}
	}
	if err := s.repo.StoreBatch(ctx, batchItems); err != nil {
		return nil, err
	}
	return shortURLS, nil
}

func (s *Shortener) HealthCheck(ctx context.Context) error {
	e := make([]error, 0)
	e = append(e, s.repo.Ping(ctx))
	return errors.Join(e...)
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
