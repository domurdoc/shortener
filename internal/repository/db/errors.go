package db

import "fmt"

type UniqueConstraintError struct {
	err error
}

func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("unique constraint violation. %s", e.err.Error())
}
