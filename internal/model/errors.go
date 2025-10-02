package model

import (
	"fmt"
	"strings"
)

type ShortCodeNotFoundError struct {
	ShortCode ShortCode
}

func (e *ShortCodeNotFoundError) Error() string {
	return fmt.Sprintf("Key %q not found", e.ShortCode)
}

type OriginalURLExistsError struct {
	OriginalURL OriginalURL
	ShortCode   ShortCode
	BatchPos    int
}

func (e *OriginalURLExistsError) Error() string {
	return fmt.Sprintf("OriginalURL %q already exists with ShortCode %q", e.OriginalURL, e.ShortCode)
}

type BatchError []*OriginalURLExistsError

func (e BatchError) Error() string {
	errorStrings := make([]string, len(e))
	for _, part := range e {
		errorStrings = append(errorStrings, part.Error())
	}
	return strings.Join(errorStrings, "\n")
}

type UserNotFoundError struct {
	UserID UserID
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User %q not found", e.UserID)
}
