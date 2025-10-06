package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/model"
)

type jsonRequest struct {
	URL string `json:"url"`
}

type jsonResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.authenticate(ctx, w, r)
	if err != nil {
		return
	}

	var req jsonRequest

	if !httputil.HasContentType(r.Header, httputil.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputil.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.Shorten(ctx, user, req.URL)
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
	h.writeJSONResponse(w, jsonResponse{Result: shortURL}, status)
}
