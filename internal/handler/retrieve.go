package handler

import (
	"errors"
	"net/http"

	"github.com/domurdoc/shortener/internal/model"
)

func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	longURL, err := h.service.GetByShortCode(r.Context(), shortCode)
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
