package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/domurdoc/shortener/internal/repository"
)

type DBRepo struct {
	db     *sql.DB
	vendor vendor
}

func New(db *sql.DB, vendor vendor) *DBRepo {
	return &DBRepo{db, vendor}
}

type vendor struct {
	queryStore string
	queryFetch string
	mapError   func(error) error
}

func (r *DBRepo) Store(ctx context.Context, key repository.Key, value repository.Value) error {
	_, err := r.db.ExecContext(
		ctx,
		r.vendor.queryStore,
		key,
		value,
	)
	var e *UniqueConstraintError
	err = r.vendor.mapError(err)
	if errors.As(err, &e) {
		return &repository.KeyAlreadyExistsError{Key: key}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *DBRepo) Fetch(ctx context.Context, key repository.Key) (repository.Value, error) {
	row := r.db.QueryRowContext(
		ctx,
		r.vendor.queryFetch,
		key,
	)
	var rawValue string
	err := row.Scan(&rawValue)
	if errors.Is(err, sql.ErrNoRows) {
		return "", &repository.KeyNotFoundError{Key: key}
	}
	if err != nil {
		return "", err
	}
	return repository.Value(rawValue), nil
}

func (r *DBRepo) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *DBRepo) Close() error {
	return r.db.Close()
}
