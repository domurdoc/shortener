package strategy

import (
	"context"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

// Strategy defines the interface for authentication strategies that handle the creation and validation of tokens.
// It abstracts the logic for generating tokens from user data and reconstructing user data from tokens.
// This allows pluggable authentication mechanisms (e.g., JWT, OAuth, session tokens).
type Strategy interface {
	// WriteToken generates a token for the given user and returns it as a string.
	// The context can be used for deadlines, cancellation, or passing request-scoped data.
	// Returns an error if token generation fails (e.g., signing error).
	WriteToken(context.Context, *model.User) (string, error)

	// ReadToken parses the given token string and retrieves the corresponding user from the repository.
	// The context can be used during user lookup (e.g., database query with timeout).
	// Returns an error if the token is invalid, expired, or user retrieval fails.
	ReadToken(context.Context, string, repository.UserRepo) (*model.User, error)
}
