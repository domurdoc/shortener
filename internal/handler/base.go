package handler

import (
	"github.com/domurdoc/shortener/internal/service"
)

type Handler struct {
	service *service.Shortener
}

func New(shortenerService *service.Shortener) *Handler {
	return &Handler{service: shortenerService}
}
