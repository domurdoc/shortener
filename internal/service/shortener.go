package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/utils"
)

// 2048 - max url length (RFC)
const (
	URLMaxLength    = 2048
	shortCodeLength = 6
)

func (s *Service) Shorten(ctx context.Context, user *model.User, originalURL string) (string, error) {
	shortCode, shortURL, err := s.generateShortCodeURL(originalURL)
	if err != nil {
		return "", err
	}
	record := &model.BaseRecord{
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
		return shortURL, urlErr
	}
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *Service) GetByShortCode(ctx context.Context, shortCode string) (string, error) {
	record, err := s.repo.Fetch(ctx, model.ShortCode(shortCode))
	if err != nil {
		return "", err
	}
	return string(record.OriginalURL), nil
}

func (s *Service) ShortenBatch(ctx context.Context, user *model.User, originalURLS []string) ([]string, error) {
	shortURLS := make([]string, 0, len(originalURLS))
	records := make([]model.BaseRecord, 0, len(originalURLS))
	for _, originalURL := range originalURLS {
		shortCode, shortURL, err := s.generateShortCodeURL(originalURL)
		if err != nil {
			return nil, err
		}
		record := model.BaseRecord{
			OriginalURL: model.OriginalURL(originalURL),
			ShortCode:   model.ShortCode(shortCode),
		}

		records = append(records, record)
		shortURLS = append(shortURLS, shortURL)
	}
	err := s.repo.StoreBatch(ctx, records, user.ID)
	var batchURLExistsErr model.BatchOriginalURLExistsError
	if errors.As(err, &batchURLExistsErr) {
		for _, urlExistsErr := range batchURLExistsErr {
			shortURL, err := url.JoinPath(s.baseURL, string(urlExistsErr.ShortCode))
			if err != nil {
				return nil, err
			}
			shortURLS[urlExistsErr.BatchPos] = shortURL
		}
		return shortURLS, batchURLExistsErr
	}
	if err != nil {
		return nil, err
	}
	return shortURLS, nil
}

func (s *Service) GetForUser(ctx context.Context, user *model.User) ([]model.URLRecord, error) {
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

func (s *Service) generateShortCodeURL(originalURL string) (string, string, error) {
	if err := validateURL(originalURL); err != nil {
		return "", "", err
	}
	shortCode := utils.GenerateRandomString(utils.ALPHA, shortCodeLength)
	shortURL, err := url.JoinPath(s.baseURL, shortCode)
	if err != nil {
		return "", "", err
	}
	return shortCode, shortURL, nil
}

func validateURL(URL string) error {
	if len(URL) > URLMaxLength {
		return &model.InvalidURLError{Msg: "url too long", URL: URL}
	}
	parsedLongURL, err := url.Parse(URL)
	if err != nil {
		return &model.InvalidURLError{Msg: err.Error(), URL: URL}
	}
	if parsedLongURL.Host == "" {
		return &model.InvalidURLError{Msg: "must be absolute", URL: URL}
	}
	if parsedLongURL.String() != URL {
		return &model.InvalidURLError{Msg: "must be url-encoded", URL: URL}
	}
	return nil
}
