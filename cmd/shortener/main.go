package main

import (
	"fmt"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/config"
)

func main() {
	options := config.LoadOptions()
	// TODO: add App struct?
	repo, log, router := app.Boot(
		options.LogLevel.String(),
		options.BaseURL.String(),
		options.DatabaseDSN.String(),
		options.FileStoragePath.String(),
	)
	defer repo.Close()
	defer log.Sync()

	log.Infow(
		"starting server",
		"addr", options.Addr,
		"baseURL", options.BaseURL,
		"logLevel", options.LogLevel,
		"fileStoragePath", options.FileStoragePath,
		"databaseDSN", options.DatabaseDSN,
		"repo", fmt.Sprintf("%T", repo),
	)
	app.Run(router, options.Addr.String(), log)
}
