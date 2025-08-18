package handler

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/domurdoc/shortener/internal/service"
)

type Shortener struct {
	baseURL string
	service *service.Shortener
}

func New(baseURL string, shortenerService *service.Shortener) *Shortener {
	return &Shortener{baseURL: baseURL, service: shortenerService}
}

func (h *Shortener) Shorten(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 2048) // 2048 - max url length (RFC)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL := string(buf[:n])
	shortCode, err := h.service.Shorten(longURL)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shortURL, err := url.JoinPath(h.baseURL, shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Shortener) GetByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	longURL, err := h.service.GetByShortCode(shortCode)
	var notFoundError *service.NotFoundError
	if errors.As(err, &notFoundError) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
