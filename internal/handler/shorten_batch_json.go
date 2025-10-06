package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/model"
)

type jsonBatchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type jsonBatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.authenticate(ctx, w, r)
	if err != nil {
		return
	}

	var reqItems []jsonBatchRequestItem

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
	shortURLS, err := h.service.ShortenBatch(ctx, user, originalURLS)
	var invalidURLErr *model.InvalidURLError
	if errors.As(err, &invalidURLErr) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var urlExistsErr model.BatchOriginalURLExistsError
	if err != nil && !errors.As(err, &urlExistsErr) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resItems := make([]jsonBatchResponseItem, len(reqItems))
	for i, jsonRequest := range reqItems {
		resItems[i] = jsonBatchResponseItem{CorrelationID: jsonRequest.CorrelationID, ShortURL: shortURLS[i]}
	}
	status := http.StatusCreated
	if err != nil {
		status = http.StatusConflict
	}
	h.writeJSONResponse(w, resItems, status)
}
