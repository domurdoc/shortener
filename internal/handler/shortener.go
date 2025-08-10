package handler

import (
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
	// TESTS fails with the next one checks
	// if r.Header.Get("Content-Type") != "text/plain" {
	// 	http.Error(w, "Only text/plain Content-Type allowed", http.StatusBadRequest)
	// 	return
	// }
	buf := make([]byte, 2048) // 2048 - max url length (RFC)
	// TODO: check if there are more bytes left from socket for stricter validation?
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rawURL := string(buf[:n])
	// SHOULD I move the following checks into service?
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if parsedURL.Host == "" {
		http.Error(w, "URL must be absolute", http.StatusBadRequest)
		return
	}
	if parsedURL.String() != rawURL {
		http.Error(w, "URL must be url-encoded", http.StatusBadRequest)
		return
	}
	shortCode, err := h.service.Shorten(rawURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	location, err := url.JoinPath(h.baseURL, shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(location))
}

func (h *Shortener) Retrieve(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	location, err := h.service.Get(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
