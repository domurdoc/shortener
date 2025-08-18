package service

import "fmt"

type URLError struct {
	msg string
	url string
}

func (e *URLError) Error() string {
	return fmt.Sprintf("Invalid URL %q: %s", e.url, e.msg)
}

type NotFoundError struct {
	shortCode string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("URL for code %q not found", e.shortCode)
}
