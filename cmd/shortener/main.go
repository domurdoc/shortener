package main

import (
	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/config"
)

func main() {
	options := config.LoadOptions()
	// TODO: add App struct?
	router, log, db := app.Boot(
		options.LogLevel.String(),
		options.FileStoragePath.String(),
		options.BaseURL.String(),
		options.DatabaseDSN.String(),
	)
	defer log.Sync()
	defer db.Close()
	app.Run(router, options.Addr.String(), log)
}
