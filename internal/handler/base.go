package handler

import (
	"database/sql"

	"github.com/domurdoc/shortener/internal/service"
)

type Shortener struct {
	service *service.Shortener
	db      *sql.DB
}

func New(shortenerService *service.Shortener, db *sql.DB) *Shortener {
	return &Shortener{service: shortenerService, db: db}
}
