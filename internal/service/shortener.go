package service

import (
	"crypto/rand"
	"net/url"

	"github.com/domurdoc/shortener/internal/repository"
)

type Shortener struct {
	repo repository.Shortener
}

const (
	charSet         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortCodeLength = 6
	charSetLength   = len(charSet)
)

func New(repo repository.Shortener) *Shortener {
	return &Shortener{repo: repo}
}

func (s *Shortener) Shorten(longURL string) (string, error) {
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
	return shortCode, s.repo.Store(repository.Key(shortCode), repository.Value(longURL))
}

func (s *Shortener) GetByShortCode(shortCode string) (string, error) {
	url, err := s.repo.Fetch(repository.Key(shortCode))
	if err != nil {
		return "", &NotFoundError{shortCode: shortCode}
	}
	return string(url), nil
}

func generateShortCode() string {
	// https://stackoverflow.com/a/67035900
	buf := make([]byte, shortCodeLength)
	rand.Read(buf)
	for i := range shortCodeLength {
		buf[i] = charSet[int(buf[i])%charSetLength]
	}
	return string(buf)
}
