package handler

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/model"
)

type jsonURLRecord struct {
	ShortURL    model.ShortURL    `json:"short_url"`
	OriginalURL model.OriginalURL `json:"original_url"`
}

// RetrieveForUser handles a request to retrieve all shortened URLs created by the authenticated user.
// It returns a JSON array of URL records or a 204 No Content if the user has no URLs.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request, which must contain user context via authentication middleware.
//
// Behavior:
//   - Extracts the authenticated user from the request context using auth.GetUser.
//   - Calls h.service.GetForUser to fetch the user's URL records.
//   - Converts the internal model.URLRecord slice to a slice of jsonURLRecord for JSON response.
//   - Responds with 200 OK if URLs exist, or 204 No Content if none.
//
// Errors:
//   - Returns 500 Internal Server Error if the service call fails.
func (h *Handler) RetrieveForUser(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	urlRecords, err := h.app.Service.GetForUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonURLRecords := make([]jsonURLRecord, 0, len(urlRecords))
	for _, ur := range urlRecords {
		jsonURLRecords = append(jsonURLRecords, jsonURLRecord(ur))
	}

	status := http.StatusOK
	if len(jsonURLRecords) == 0 {
		status = http.StatusNoContent
	}

	h.writeJSONResponse(w, jsonURLRecords, status)
}
