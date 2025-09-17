package db

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgQueryStore = `INSERT INTO records (key, value) VALUES ($1, $2)`
	pgQueryFetch = `SELECT value FROM records WHERE key = $1`
)

func pgMapError(err error) error {
	var pe *pgconn.PgError
	if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
		return &UniqueConstraintError{pe}
	}
	return err
}

func NewPGVendor() *vendor {
	return &vendor{
		queryStore: pgQueryStore,
		queryFetch: pgQueryFetch,
		mapError:   pgMapError,
	}
}
