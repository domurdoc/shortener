package handler

import (
	"github.com/domurdoc/shortener/internal/service"
)

type Shortener struct {
	service *service.Shortener
}

func New(shortenerService *service.Shortener) *Shortener {
	return &Shortener{service: shortenerService}
}
