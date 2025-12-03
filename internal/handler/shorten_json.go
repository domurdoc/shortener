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

type jsonRequest struct {
	URL string `json:"url"`
}

type jsonResponse struct {
	Result string `json:"result"`
}

// ShortenJSON handles a JSON-based request to shorten a single URL.
// It expects a JSON object with a "url" field and returns a JSON object with the shortened URL.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request containing the JSON body and user context.
//
// Behavior:
//   - Extracts the authenticated user from the request context.
//   - Validates that the Content-Type is application/json.
//   - Decodes the request body into a jsonRequest.
//   - Calls h.service.Shorten to generate the short URL.
//   - Records the shortening action in the audit log.
//
// Status Codes:
//   - 201 Created on success.
//   - 400 Bad Request if Content-Type is invalid, JSON is malformed, or URL is invalid.
//   - 409 Conflict if the URL already exists (with existing short URL in response).
//   - 500 Internal Server Error for unexpected service errors.
//
// Error Handling:
//   - *model.InvalidURLError results in 400.
//   - *model.OriginalURLExistsError results in 409 with the existing short URL.
func (h *Handler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	var req jsonRequest

	user := auth.GetUser(r)

	if !httputil.HasContentType(r.Header, httputil.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputil.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := h.app.Service.Shorten(r.Context(), user, req.URL)
	h.app.Audit.Shortened(user.ID, req.URL)
	var invalidURLErr *model.InvalidURLError
	if errors.As(err, &invalidURLErr) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var urlExistsErr *model.OriginalURLExistsError
	if err != nil && !errors.As(err, &urlExistsErr) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	status := http.StatusCreated
	if err != nil {
		status = http.StatusConflict
	}
	h.writeJSONResponse(w, jsonResponse{Result: shortURL}, status)
}
