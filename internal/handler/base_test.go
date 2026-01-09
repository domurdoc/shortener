package handler

import (
	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/config"
)

func getAppHandler() (*app.App, *Handler) {
	cfg := config.Default()
	app, err := app.New(cfg)
	if err != nil {
		panic(err)
	}
	return app, New(app)
}
