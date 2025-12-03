package utils

import (
	"errors"
	"slices"
)

// Closer is a utility type that manages multiple closeable resources.
// It allows registering multiple close functions and invoking them all at once,
// aggregating any errors returned by individual closers.
// This is useful for cleaning up resources like database connections, file handles, or network listeners.
type Closer struct {
	closes []func() error // closes holds a list of functions to be called during Close.
}

func NewCloser() *Closer {
	return &Closer{}
}

func (c *Closer) Register(close func() error) {
	c.closes = append(c.closes, close)
}

func (c *Closer) Close() error {
	errs := make([]error, 0, len(c.closes))
	slices.Reverse(c.closes)
	for _, c := range c.closes {
		errs = append(errs, c())
	}
	c.closes = nil
	return errors.Join(errs...)
}
