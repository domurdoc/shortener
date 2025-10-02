package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/service"
)

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.authenticate(ctx, w, r)
	if err != nil {
		return
	}

	buf := make([]byte, service.URLMaxLength)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL := string(buf[:n])
	shortURL, err := h.service.Shorten(ctx, user, longURL)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && !errors.Is(err, service.ErrURLConflict) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	status := http.StatusCreated
	if err != nil {
		status = http.StatusConflict
	}
	httputil.SetContentType(w.Header(), httputil.ContentTypeTextPlain)
	w.WriteHeader(status)
	w.Write([]byte(shortURL))
}
