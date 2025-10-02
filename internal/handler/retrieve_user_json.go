package handler

import (
	"net/http"

	"github.com/domurdoc/shortener/internal/model"
)

type jsonURLRecord struct {
	ShortURL    model.ShortURL    `json:"short_url"`
	OriginalURL model.OriginalURL `json:"original_url"`
}

func (h *Handler) RetrieveForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.authenticate(ctx, w, r)
	if err != nil {
		return
	}

	urlRecords, err := h.service.GetForUser(ctx, user)
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
