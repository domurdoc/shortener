package strategy

import (
	"context"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

type Strategy interface {
	WriteToken(context.Context, *model.User) (string, error)
	ReadToken(context.Context, string, repository.UserRepo) (*model.User, error)
}
