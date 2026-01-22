package handler

import (
	"context"
	"net/http"
	"time"
)

type jsonStats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// GetStats returns the number of URLs and users
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	urlsCount, urlsErr := h.app.Repos.Record.CountURLs(ctx)
	if urlsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	usersCount, usersErr := h.app.Repos.User.CountUsers(ctx)
	if usersErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	stats := jsonStats{Urls: urlsCount, Users: usersCount}
	h.writeJSONResponse(w, stats, http.StatusOK)
}
