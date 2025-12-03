package handler

import (
	"encoding/json"
	"net/http"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/httputil"
)

// Handler encapsulates HTTP request handlers for the shortener application.
// It acts as a bridge between HTTP routes and business logic in the service layer.
// It also includes an audit component for tracking user actions.
type Handler struct {
	app *app.App
}

func New(a *app.App) *Handler {
	return &Handler{app: a}
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
