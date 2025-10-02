package main

import (
	"fmt"
	"log"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/router"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	a.Log.Infow(
		"starting server",
		"addr", a.Options.Addr,
		"baseURL", a.Options.BaseURL,
		"logLevel", a.Options.LogLevel,
		"fileStoragePath", a.Options.FileStoragePath,
		"databaseDSN", a.Options.DatabaseDSN,
		"repo", fmt.Sprintf("%T", a.RecordRepo),
	)
	log.Fatal(router.Run(a.Service, a.Auth, a.DB, a.Log, a.Options.Addr.String()))
}
