package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/domurdoc/shortener/migrations"
)

func NewPG(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigratePG(pgDB *sql.DB) error {
	d1, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}
	d2, err := postgres.WithInstance(pgDB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", d1, "postgres", d2)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func NewPGArger() Arger {
	return &pgArger{}
}

type pgArger struct {
	Pos int
}

func (a *pgArger) Next() string {
	a.Pos += 1
	return fmt.Sprintf("$%d", a.Pos)
}
