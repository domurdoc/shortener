package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/service"
)

type jsonRequest struct {
	URL string `json:"url"`
}

type jsonResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	var req jsonRequest

	if !httputil.HasContentType(r.Header, httputil.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputil.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.Shorten(r.Context(), req.URL)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, service.ErrURLConflict) {
		writeShortURLJSON(w, http.StatusConflict, shortURL)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeShortURLJSON(w, http.StatusCreated, shortURL)
}

func writeShortURLJSON(w http.ResponseWriter, status int, shortURL string) {
	httputil.SetContentType(w.Header(), httputil.ContentTypeJSON)
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(jsonResponse{Result: shortURL}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
