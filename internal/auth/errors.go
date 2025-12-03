package auth

import "fmt"

// NoTokenError represents an error that occurs when an authentication token is expected but not found.
// It wraps the underlying error that led to the detection of a missing token.
// This error is typically returned by authentication strategies when parsing or retrieving a token fails
// due to its absence in the request (e.g., missing cookie or authorization header).
type NoTokenError struct {
	Err error // Err is the underlying error that caused or describes the absence of a token.
}

func (e *NoTokenError) Error() string {
	return fmt.Sprintf("no token: %v", e.Err)

}

func (e *NoTokenError) Unwrap() error {
	return e.Err
}

type InvalidTokenError struct {
	Err error
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("invalid token: %v", e.Err)

}

func (e *InvalidTokenError) Unwrap() error {
	return e.Err
}
