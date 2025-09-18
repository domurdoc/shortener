package repository

import (
	"fmt"
	"strings"
)

type KeyNotFoundError struct {
	Key Key
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key %q not found", e.Key)
}

type KeyAlreadyExistsError struct {
	Key Key
}

func (e *KeyAlreadyExistsError) Error() string {
	return fmt.Sprintf("Key %q already exists", e.Key)
}

type ValueAlreadyExistsError struct {
	Key   Key
	Value Value
	Pos   int
}

func (e *ValueAlreadyExistsError) Error() string {
	return fmt.Sprintf("Value %q already exists with key %q", e.Value, e.Key)
}

type BatchError []*ValueAlreadyExistsError

func (e BatchError) Error() string {
	errorStrings := make([]string, len(e))
	for _, part := range e {
		errorStrings = append(errorStrings, part.Error())
	}
	return strings.Join(errorStrings, "\n")
}
