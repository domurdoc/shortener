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

// Auth is responsible for handling user authentication and authorization in the application.
// It combines a strategy for token management (e.g., JWT) and a transport mechanism (e.g., cookies)
// to authenticate requests and manage user sessions.
// It also interacts with a user repository to load user data when needed.
type Auth struct {
	strategy  strategy.Strategy   // strategy defines how tokens are created, parsed, and validated (e.g., JWT).
	transport transport.Transport // transport defines how authentication data is exchanged with the client (e.g., via cookies).
	repo      repository.UserRepo // repo is used to retrieve user information during authentication when necessary.
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
