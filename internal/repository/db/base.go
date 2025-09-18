package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/domurdoc/shortener/internal/repository"
)

type DBRepo struct {
	db       *sql.DB
	newArger func() arger
}

func New(db *sql.DB, newArger func() arger) *DBRepo {
	return &DBRepo{db, newArger}
}

type arger interface {
	next() string
}

const (
	queryStore = `
INSERT INTO records (key, value) VALUES (%s, %s)
ON CONFLICT (value) DO UPDATE SET key = records.key
RETURNING key
`
	queryFetch = `
SELECT value FROM records WHERE key = %s
`
)

func (r *DBRepo) Store(ctx context.Context, key repository.Key, value repository.Value) error {
	arger := r.newArger()
	query := fmt.Sprintf(queryStore, arger.next(), arger.next())

	row := r.db.QueryRowContext(ctx, query, key, value)
	var returnedKey repository.Key
	err := row.Scan(&returnedKey)
	if err != nil {
		return err
	}
	if returnedKey != key {
		return &repository.ValueAlreadyExistsError{Key: returnedKey, Value: value}
	}
	return nil
}

func (r *DBRepo) StoreBatch(ctx context.Context, batchItems []repository.BatchItem) error {
	arger := r.newArger()
	query := fmt.Sprintf(queryStore, arger.next(), arger.next())

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var batchError repository.BatchError
	for pos, item := range batchItems {
		row := stmt.QueryRowContext(ctx, item.Key, item.Value)
		var returnedKey repository.Key
		err := row.Scan(&returnedKey)
		if err != nil {
			return err
		}
		if returnedKey != item.Key {
			valueErr := &repository.ValueAlreadyExistsError{Key: returnedKey, Value: item.Value, Pos: pos}
			batchError = append(batchError, valueErr)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	if len(batchError) != 0 {
		return batchError
	}
	return nil
}

func (r *DBRepo) Fetch(ctx context.Context, key repository.Key) (repository.Value, error) {
	var value repository.Value

	arger := r.newArger()
	query := fmt.Sprintf(queryFetch, arger.next())

	row := r.db.QueryRowContext(ctx, query, key)
	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", &repository.KeyNotFoundError{Key: key}
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *DBRepo) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *DBRepo) Close() error {
	return r.db.Close()
}
