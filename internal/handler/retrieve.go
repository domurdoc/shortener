package handler

import (
	"errors"
	"net/http"

	"github.com/domurdoc/shortener/internal/model"
)

// Retrieve handles redirect requests for a given short code.
// It looks up the original URL and issues a 307 Temporary Redirect.
// If the short code is not found or has been deleted, it returns appropriate error statuses.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request, which should contain the shortCode as a path parameter.
//
// Behavior:
//   - Extracts the shortCode from the request path using r.PathValue.
//   - Calls h.service.GetByShortCode to retrieve the original URL.
//   - Records the access in the audit log via h.audit.Followed (with zero user ID for anonymous access).
//   - On success, responds with 307 Temporary Redirect and sets the Location header.
//
// Error Handling:
//   - 404 Not Found if the short code does not exist (*model.ShortCodeNotFoundError).
//   - 410 Gone if the short code exists but was deleted (*model.ShortCodeDeletedError).
//   - 500 Internal Server Error for any other error.
func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	longURL, err := h.app.Service.GetByShortCode(r.Context(), shortCode)
	h.app.Audit.Audit.Followed(model.UserID(0), longURL)
	var notFoundErr *model.ShortCodeNotFoundError
	if errors.As(err, &notFoundErr) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var isDeletedErr *model.ShortCodeDeletedError
	if errors.As(err, &isDeletedErr) {
		http.Error(w, "", http.StatusGone)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
