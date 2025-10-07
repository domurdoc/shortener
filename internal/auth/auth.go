package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/domurdoc/shortener/internal/auth/strategy"
	"github.com/domurdoc/shortener/internal/auth/transport"
	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

type Auth struct {
	strategy  strategy.Strategy
	transport transport.Transport
	repo      repository.UserRepo
}

func New(strategy strategy.Strategy, transport transport.Transport, repo repository.UserRepo) *Auth {
	return &Auth{
		strategy:  strategy,
		transport: transport,
		repo:      repo,
	}
}

func (a *Auth) Authenticate(ctx context.Context, r *http.Request) (*model.User, error) {
	tokenString, err := a.transport.Read(r)
	if err != nil {
		return nil, &NoTokenError{err}
	}
	user, err := a.strategy.ReadToken(ctx, tokenString, a.repo)
	if err != nil {
		return nil, &InvalidTokenError{err}
	}
	return user, nil
}

func (a *Auth) Login(ctx context.Context, w http.ResponseWriter, user *model.User) error {
	tokenString, err := a.strategy.WriteToken(ctx, user)
	if err != nil {
		return err
	}
	return a.transport.Write(w, tokenString)
}

func (a *Auth) Register(ctx context.Context) (*model.User, error) {
	return a.repo.CreateUser(ctx)
}

func (a *Auth) AuthenticateOrRegisterAndLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	user, err := a.Authenticate(ctx, r)
	if err != nil {
		var noTokenErr *NoTokenError
		if errors.As(err, &noTokenErr) {
			user, err = a.Register(ctx)
			if err != nil {
				return nil, err
			}
			if err = a.Login(ctx, w, user); err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	return user, nil
}
