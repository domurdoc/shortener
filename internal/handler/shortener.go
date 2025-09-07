package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputils"
	"github.com/domurdoc/shortener/internal/service"
)

type Shortener struct {
	service *service.Shortener
}

func New(shortenerService *service.Shortener) *Shortener {
	return &Shortener{service: shortenerService}
}

func (h *Shortener) Shorten(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 2048) // 2048 - max url length (RFC)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL := string(buf[:n])
	shortURL, err := h.service.Shorten(longURL)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputils.SetContentType(w.Header(), httputils.ContentTypeTextPlain)
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

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func (h *Shortener) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	var req Request

	if !httputils.HasContentType(r.Header, httputils.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputils.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.Shorten(req.URL)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputils.SetContentType(w.Header(), httputils.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(Response{Result: shortURL}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
