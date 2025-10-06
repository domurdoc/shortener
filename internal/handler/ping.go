package handler

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := h.service.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
