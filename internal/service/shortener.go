package service

import (
	"fmt"
	"math/rand"

	"github.com/domurdoc/shortener/internal/repository"
)

type Shortener struct {
	repo    repository.Shortener
	genCode func() string
}

const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const defaultLength = 6

func New(repo repository.Shortener, genCode func() string) *Shortener {
	if genCode == nil {
		genCode = NewGenFunc(defaultCharset, defaultLength)
	}
	return &Shortener{
		repo:    repo,
		genCode: genCode,
	}
}

type NotFoundError struct {
	shortCode string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("URL for code %q not found", e.shortCode)
}

func (s *Shortener) Shorten(URL string) (string, error) {
	// URL expected to be valid absolute url-encoded
	code := s.genCode()
	return code, s.repo.Store(repository.Key(code), repository.Value(URL))
}

func (s *Shortener) Get(shortCode string) (string, error) {
	url, err := s.repo.Fetch(repository.Key(shortCode))
	if err != nil {
		return "", &NotFoundError{shortCode: shortCode}
	}
	return string(url), nil
}

func NewGenFunc(charset string, length int) func() string {
	return func() string {
		chars := make([]byte, length)
		for i := range length {
			chars[i] = charset[rand.Intn(len(charset))]
		}
		return string(chars)
	}
}
