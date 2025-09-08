package main

import (
	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/config"
)

func main() {
	options := config.LoadOptions()
	router, log := app.Boot(
		options.LogLevel.String(),
		options.FileStoragePath.String(),
		options.BaseURL.String(),
	)
	defer log.Sync()
	app.Run(router, options.Addr.String(), log)
}
