package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/httputil"
)

// DeleteShortCodes handles a request to delete multiple shortened URLs for the authenticated user.
// It processes the deletion asynchronously, returning immediately with a 202 Accepted status.
//
// The request must:
//   - Use Content-Type: application/json
//   - Contain a JSON array of short codes (strings) in the request body
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request, which should contain the user context and JSON body.
//
// Behavior:
//   - Extracts the authenticated user from the request context via auth.GetUser.
//   - Validates the Content-Type header.
//   - Decodes the JSON array of short codes.
//   - Initiates asynchronous deletion via h.service.DeleteShortCodes in a goroutine.
//   - Responds with 202 Accepted upon successful initiation.
//
// Errors:
//   - Returns 400 Bad Request if Content-Type is incorrect or JSON is malformed.
//
// Note: The actual deletion is handled asynchronously; failure during deletion
// is not reported back to the client.
func (h *Handler) DeleteShortCodes(w http.ResponseWriter, r *http.Request) {
	var shortCodes []string

	user := auth.GetUser(r)

	if !httputil.HasContentType(r.Header, httputil.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputil.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&shortCodes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	go h.app.Service.DeleteShortCodes(r.Context(), user, shortCodes)
	httputil.SetContentType(w.Header(), httputil.ContentTypeJSON)
	w.WriteHeader(http.StatusAccepted)
}
