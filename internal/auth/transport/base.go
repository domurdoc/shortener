package transport

import "net/http"

type Transport interface {
	Read(*http.Request) (string, error)
	Write(http.ResponseWriter, string) error
}
