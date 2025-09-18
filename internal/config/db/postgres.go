package db

import (
	"database/sql"
	"errors"

	"github.com/domurdoc/shortener"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPG(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func MigratePG(pgDB *sql.DB) {
	d1, err := iofs.New(shortener.FS, "migrations")
	if err != nil {
		panic(err)
	}
	d2, err := postgres.WithInstance(pgDB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithInstance("iofs", d1, "postgres", d2)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
