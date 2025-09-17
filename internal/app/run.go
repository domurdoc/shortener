package app

import (
	"net/http"

	"go.uber.org/zap"
)

func Run(router http.Handler, address string, log *zap.SugaredLogger) {
	if err := http.ListenAndServe(address, router); err != nil {
		log.Errorw(
			err.Error(),
			"event", "start server",
		)
		panic(err)
	}
}
