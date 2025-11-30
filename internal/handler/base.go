package handler

import (
	"encoding/json"
	"net/http"

	"github.com/domurdoc/shortener/internal/audit"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/service"
)

type Handler struct {
	service *service.Service
	audit   *audit.Audit
}

func New(service *service.Service, audit *audit.Audit) *Handler {
	return &Handler{service: service, audit: audit}
}

func (h *Handler) writeJSONResponse(w http.ResponseWriter, response any, status int) {
	httputil.SetContentType(w.Header(), httputil.ContentTypeJSON)
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
