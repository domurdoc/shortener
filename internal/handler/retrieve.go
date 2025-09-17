package handler

import (
	"errors"
	"net/http"

	"github.com/domurdoc/shortener/internal/service"
)

func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	longURL, err := h.service.GetByShortCode(r.Context(), shortCode)
	var notFoundError *service.NotFoundError
	if errors.As(err, &notFoundError) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
