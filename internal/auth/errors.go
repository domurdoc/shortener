package auth

import "fmt"

type NoTokenError struct {
	Err error
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
