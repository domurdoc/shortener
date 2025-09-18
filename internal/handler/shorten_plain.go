package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/service"
)

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, service.URLMaxLength)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL := string(buf[:n])
	shortURL, err := h.service.Shorten(r.Context(), longURL)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, service.ErrURLConflict) {
		writeShortURL(w, http.StatusConflict, shortURL)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeShortURL(w, http.StatusCreated, shortURL)
}

func writeShortURL(w http.ResponseWriter, status int, shortURL string) {
	httputil.SetContentType(w.Header(), httputil.ContentTypeTextPlain)
	w.WriteHeader(status)
	w.Write([]byte(shortURL))
}
