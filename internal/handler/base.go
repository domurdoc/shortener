package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/service"
)

type Handler struct {
	service *service.Shortener
	auth    *auth.Auth
	db      *sql.DB
}

func New(shortenerService *service.Shortener, auth *auth.Auth, db *sql.DB) *Handler {
	return &Handler{service: shortenerService, auth: auth, db: db}
}

func (h *Handler) authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	user, err := h.auth.Authenticate(ctx, r)

	var noTokenErr *auth.NoTokenError
	if errors.As(err, &noTokenErr) {
		user, err := h.auth.Register(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		if err = h.auth.Login(ctx, w, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		return user, nil
	}

	var invalidTokenErr *auth.InvalidTokenError
	if errors.As(err, &invalidTokenErr) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	}

	return user, nil
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
