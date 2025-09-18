package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/service"
)

type jsonBatchRequestItem struct {
	CID         string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type jsonBatchResponseItem struct {
	CID      string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func (h *Handler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
	var reqItems []jsonBatchRequestItem
	ctx := r.Context()

	if !httputil.HasContentType(r.Header, httputil.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputil.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqItems); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(reqItems) == 0 {
		http.Error(w, "at least one item must be passed", http.StatusBadRequest)
		return
	}
	originalURLS := make([]string, len(reqItems))
	for i, jsonRequest := range reqItems {
		originalURLS[i] = jsonRequest.OriginalURL
	}

	shortURLS, err := h.service.ShortenBatch(ctx, originalURLS)
	var urlError *service.URLError
	if errors.As(err, &urlError) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, service.ErrURLConflict) {
		writeShortURLBatchJSON(w, http.StatusConflict, shortURLS, reqItems)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeShortURLBatchJSON(w, http.StatusCreated, shortURLS, reqItems)
}

func writeShortURLBatchJSON(w http.ResponseWriter, status int, shortURLS []string, reqItems []jsonBatchRequestItem) {
	resItems := make([]jsonBatchResponseItem, len(reqItems))
	for i, jsonRequest := range reqItems {
		resItems[i] = jsonBatchResponseItem{CID: jsonRequest.CID, ShortURL: shortURLS[i]}
	}

	httputil.SetContentType(w.Header(), httputil.ContentTypeJSON)
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
