package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log = zap.NewNop()
var Sugar zap.SugaredLogger

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	Sugar = *zl.Sugar()
	return nil
}

type (
	requestData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		requestData *requestData
	}
)

func (w loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.requestData.size += n
	return n, err
}

func (w loggingResponseWriter) WriteHeader(statusCode int) {
	w.requestData.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		requestData := requestData{}
		w = loggingResponseWriter{w, &requestData}
		h.ServeHTTP(w, r)
		duration := time.Since(start)
		Sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", requestData.status,
			"size", requestData.size,
		)
	})
}
