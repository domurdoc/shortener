package repository

import (
	"context"

	"github.com/domurdoc/shortener/internal/model"
)

type RecordRepo interface {
	Store(context.Context, *model.BaseRecord, model.UserID) error
	Fetch(context.Context, model.ShortCode) (*model.BaseRecord, error)
	FetchForUser(context.Context, model.UserID) ([]model.BaseRecord, error)
	StoreBatch(context.Context, []model.BaseRecord, model.UserID) error
	Delete(context.Context, []model.UserRecord) (int, error)
}

type UserRepo interface {
	GetUser(context.Context, model.UserID) (*model.User, error)
	CreateUser(context.Context) (*model.User, error)
}
