package repository

import (
	"context"

	"github.com/domurdoc/shortener/internal/model"
)

type RecordRepo interface {
	Store(context.Context, *model.Record, model.UserID) error
	Fetch(context.Context, model.ShortCode) (*model.Record, error)
	FetchForUser(context.Context, model.UserID) ([]model.Record, error)
	StoreBatch(context.Context, []model.Record, model.UserID) error
}

type UserRepo interface {
	GetUser(context.Context, model.UserID) (*model.User, error)
	CreateUser(context.Context) (*model.User, error)
}
