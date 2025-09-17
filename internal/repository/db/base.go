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
	return r.StoreBatch(ctx, repository.SingleItemBatch(key, value))
}

func (r *DBRepo) StoreBatch(ctx context.Context, batchItems []repository.BatchItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, r.vendor.queryStore)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range batchItems {
		_, err := stmt.ExecContext(ctx, item.Key, item.Value)
		var e *UniqueConstraintError
		err = r.vendor.mapError(err)
		if errors.As(err, &e) {
			return &repository.KeyAlreadyExistsError{Key: item.Key}
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *DBRepo) Fetch(ctx context.Context, key repository.Key) (repository.Value, error) {
	row := r.db.QueryRowContext(ctx, r.vendor.queryFetch, key)
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
