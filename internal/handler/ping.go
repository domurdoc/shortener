package handler

import (
	"context"
	"net/http"
	"time"
)

// Ping handles health check requests by verifying the database connection.
// It uses a 1-second timeout to prevent long delays during outages.
// This endpoint is typically used by load balancers or monitoring tools.
//
// If the database responds within the timeout, it returns 200 OK.
// If the database is unreachable or the request times out, it returns 500 Internal Server Error.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
//
// The context from the incoming request is wrapped with a timeout and passed to the service layer.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := h.app.Service.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
