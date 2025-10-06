package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
)

func (h *Handler) DeleteShortCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.authenticate(ctx, w, r)
	if err != nil {
		return
	}
	var shortCodes []string
	if !httputil.HasContentType(r.Header, httputil.ContentTypeJSON) {
		http.Error(w, fmt.Sprintf("wanted Content-Type: %s", httputil.ContentTypeJSON), http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&shortCodes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	go h.service.DeleteShortCodes(ctx, user, shortCodes)
	httputil.SetContentType(w.Header(), httputil.ContentTypeJSON)
	w.WriteHeader(http.StatusAccepted)
}
