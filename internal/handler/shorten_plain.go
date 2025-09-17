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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputil.SetContentType(w.Header(), httputil.ContentTypeTextPlain)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
