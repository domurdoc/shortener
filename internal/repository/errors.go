package repository

import "fmt"

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
