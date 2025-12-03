package model

import (
	"fmt"
	"strings"
)

// InvalidURLError represents an error when a provided URL fails validation.
// This may be due to invalid format, missing host, or exceeding length limits.
type InvalidURLError struct {
	// Msg describes the specific validation failure.
	Msg string
	// URL is the invalid URL string that caused the error.
	URL string
}

// Error returns the string representation of the InvalidURLError.
func (e InvalidURLError) Error() string {
	return fmt.Sprintf("Invalid URL %q: %s", e.URL, e.Msg)
}

// ShortCodeNotFoundError is returned when a requested short code does not exist in the system.
type ShortCodeNotFoundError struct {
	// ShortCode is the non-existent short code that was requested.
	ShortCode ShortCode
}

// Error returns the string representation of the ShortCodeNotFoundError.
func (e ShortCodeNotFoundError) Error() string {
	return fmt.Sprintf("ShortCode %q not found", e.ShortCode)
}

// ShortCodeDeletedError is returned when a short code exists but has been marked as deleted (no active owners).
type ShortCodeDeletedError struct {
	// ShortCode is the deleted short code that was requested.
	ShortCode ShortCode
}

// Error returns the string representation of the ShortCodeDeletedError.
func (e ShortCodeDeletedError) Error() string {
	return fmt.Sprintf("ShortCode %q deleted", e.ShortCode)
}

// OriginalURLExistsError is returned when attempting to shorten a URL that already exists for a user.
// It includes the existing short code to allow clients to retrieve the existing short URL.
type OriginalURLExistsError struct {
	// OriginalURL is the URL that already exists in the system.
	OriginalURL OriginalURL
	// ShortCode is the existing short code for the URL.
	ShortCode ShortCode
	// BatchPos indicates the position of the URL in a batch request, if applicable.
	BatchPos int
}

// Error returns the string representation of the OriginalURLExistsError.
func (e OriginalURLExistsError) Error() string {
	return fmt.Sprintf("OriginalURL %q already exists with ShortCode %q", e.OriginalURL, e.ShortCode)
}

// BatchOriginalURLExistsError is a slice of OriginalURLExistsError instances.
// It is returned during batch operations when multiple URLs already exist.
// Each entry contains details about a duplicate URL and its position in the batch.
type BatchOriginalURLExistsError []*OriginalURLExistsError

// Error returns the string representation of the BatchOriginalURLExistsError,
// concatenating all individual errors with newlines.
func (e BatchOriginalURLExistsError) Error() string {
	errorStrings := make([]string, len(e))
	for _, part := range e {
		errorStrings = append(errorStrings, part.Error())
	}
	return strings.Join(errorStrings, "\n")
}

// UserNotFoundError is returned when a requested user ID does not exist in the system.
type UserNotFoundError struct {
	// UserID is the identifier of the non-existent user.
	UserID UserID
}

// Error returns the string representation of the UserNotFoundError.
func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User %q not found", e.UserID)
}
