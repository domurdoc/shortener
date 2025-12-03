package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/auth"
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

// ShortenBatchJSON handles a JSON-based batch request to shorten multiple URLs at once.
// It preserves the order of requests and includes correlation IDs in the response.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request, which must contain a JSON array of jsonBatchRequestItem.
//
// Behavior:
//   - Extracts the authenticated user from the request context.
//   - Validates that the Content-Type is application/json.
//   - Decodes the JSON request body into a slice of jsonBatchRequestItem.
//   - Ensures at least one item is provided.
//   - Extracts original URLs and passes them to h.service.ShortenBatch.
//   - Constructs a response with correlation IDs preserved.
//
// Status Codes:
//   - 201 Created on full success.
//   - 400 Bad Request if Content-Type is invalid, JSON is malformed, or no items are provided.
//   - 409 Conflict if some URLs already exist (partial success with conflict details).
//   - 500 Internal Server Error for unexpected service errors.
//
// Error Handling:
//   - *model.InvalidURLError results in 400.
//   - model.BatchOriginalURLExistsError results in 409 with partial response.
func (h *Handler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
	var reqItems []jsonBatchRequestItem

	user := auth.GetUser(r)

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
	shortURLS, err := h.app.Service.ShortenBatch(r.Context(), user, originalURLS)
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
