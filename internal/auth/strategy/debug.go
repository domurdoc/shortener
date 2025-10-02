package strategy

import (
	"context"
	"strconv"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

type DebugStrategy struct{}

func NewDebug() *DebugStrategy {
	return &DebugStrategy{}
}

func (s *DebugStrategy) WriteToken(ctx context.Context, user *model.User) (string, error) {
	return strconv.Itoa(int(user.ID)), nil
}

func (s *DebugStrategy) ReadToken(ctx context.Context, tokenString string, repo repository.UserRepo) (*model.User, error) {
	userID, err := strconv.Atoi(tokenString)
	if err != nil {
		return nil, err
	}
	return repo.GetUser(ctx, model.UserID(userID))
}
