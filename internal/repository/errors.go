package repository

import "fmt"

type KeyNotFoundError struct {
	key Key
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key %q not found", e.key)
}

type KeyAlreadyExistsError struct {
	key Key
}

func (e *KeyAlreadyExistsError) Error() string {
	return fmt.Sprintf("Key %q already exists", e.key)
}
