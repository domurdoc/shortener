package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/service"
)

// Shorten handles a form-based or raw POST request to shorten a URL.
// The original URL is read directly from the request body.
// It returns the short URL as plain text.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request with the original URL in the body and user context.
//
// Behavior:
//   - Extracts the authenticated user from the request context.
//   - Reads up to service.URLMaxLength bytes from the request body.
//   - Calls h.service.Shorten to generate the short URL.
//   - Records the shortening action in the audit log.
//
// Response:
//   - 201 Created on success, with the short URL in the response body (text/plain).
//   - 400 Bad Request if the URL is invalid or reading fails.
//   - 409 Conflict if the URL already exists (with existing short URL in body).
//   - 500 Internal Server Error for unexpected service errors.
//
// The response Content-Type is set to "text/plain; charset=utf-8".
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	buf := make([]byte, service.URLMaxLength)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL := string(buf[:n])
	shortURL, err := h.app.Service.Shorten(r.Context(), user, longURL)
	h.app.Audit.Audit.Shortened(user.ID, longURL)
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
	httputil.SetContentType(w.Header(), httputil.ContentTypeTextPlain)
	w.WriteHeader(status)
	w.Write([]byte(shortURL))
}
