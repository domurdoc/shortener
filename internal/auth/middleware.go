package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/model"
)

type ctxKey string

const userKey = ctxKey("user")

func NewAuthMiddleware(auth *Auth) httputil.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := auth.AuthenticateOrRegisterAndLogin(ctx, w, r)
			if err != nil {
				var invalidTokenErr *InvalidTokenError
				if errors.As(err, &invalidTokenErr) {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			h.ServeHTTP(w, AttachUser(r, user))
		})
	}
}

func GetUser(r *http.Request) *model.User {
	return r.Context().Value(userKey).(*model.User)
}

func AttachUser(r *http.Request, user *model.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userKey, user))
}
