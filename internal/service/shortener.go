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
	if len(longURL) > URLMaxLength {
		return "", &URLError{msg: "url too long", url: longURL}
	}
	parsedLongURL, err := url.Parse(longURL)
	if err != nil {
		return "", &URLError{msg: err.Error(), url: longURL}
	}
	if parsedLongURL.Host == "" {
		return "", &URLError{msg: "must be absolute", url: longURL}
	}
	if parsedLongURL.String() != longURL {
		return "", &URLError{msg: "must be url-encoded", url: longURL}
	}
	shortCode := generateShortCode()
	shortURL, err := url.JoinPath(s.baseURL, shortCode)
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

func (s *Shortener) HealthCheck(ctx context.Context) error {
	e := make([]error, 0)
	e = append(e, s.repo.Ping(ctx))
	return errors.Join(e...)
}
