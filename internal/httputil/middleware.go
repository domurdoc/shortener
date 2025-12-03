package httputil

import "net/http"

// Middleware is a function type that takes an http.Handler and returns another http.Handler.
// It is used to compose layered functionality (e.g., logging, authentication, compression)
// around HTTP request handlers.
//
// This type alias allows for reusable and chainable middleware components in HTTP routing.
// Each middleware can inspect or modify the request and response, or terminate the chain early.
//
// Example implementations include:
//   - Logging middleware to record request details
//   - Authentication middleware to validate user sessions
//   - Compression middleware to handle GZIP encoding
//
// Middlewares are typically combined using helper functions like AddMiddlewares
// to wrap a final handler with multiple cross-cutting concerns.
type Middleware func(http.Handler) http.Handler

func AddMiddlewares(root http.Handler, middlewares ...Middleware) http.Handler {
	handler := root
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	return handler
}
