package app

import (
	"net/http"

	"go.uber.org/zap"
)

func Run(router http.Handler, address string, log *zap.SugaredLogger) {
	log.Infow(
		"Starting server",
		"addr", address,
	)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Errorw(
			err.Error(),
			"event", "start server",
		)
		panic(err)
	}
}
