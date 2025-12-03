package transport

import "net/http"

// Transport defines how authentication data (such as tokens) is transmitted between the client and server.
// It abstracts the mechanism for reading tokens from incoming requests and writing them to outgoing responses.
// Common implementations may use HTTP headers, cookies, or other transport methods.
type Transport interface {
	// Read extracts the authentication token from the given HTTP request.
	// Returns the token as a string on success, or an error if the token is missing or malformed.
	Read(*http.Request) (string, error)

	// Write sends the authentication token to the client via the HTTP response writer.
	// The token is usually embedded in cookies or headers.
	// Returns an error if writing fails (e.g., response already written).
	Write(http.ResponseWriter, string) error
}
