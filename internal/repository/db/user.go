package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/model"
)

type DBUserRepo struct {
	db       *sql.DB
	newArger func() db.Arger
}

func NewDBUserRepo(db *sql.DB, newArger func() db.Arger) *DBUserRepo {
	return &DBUserRepo{db, newArger}
}

const (
	queryCreateUser = `
INSERT INTO users DEFAULT VALUES RETURNING id
`
	queryGetUser = `
SELECT id FROM users WHERE id = %s
`
)

func (r *DBUserRepo) CreateUser(ctx context.Context) (*model.User, error) {
	var user model.User

	row := r.db.QueryRowContext(ctx, queryCreateUser)
	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *DBUserRepo) GetUser(ctx context.Context, userID model.UserID) (*model.User, error) {
	var user model.User

	arger := r.newArger()
	query := fmt.Sprintf(queryGetUser, arger.Next())

	row := r.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &model.UserNotFoundError{UserID: userID}
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
