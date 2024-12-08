package errors

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("record with ID %s not found", e.ID)
}

func (e *NotFoundError) Is(target error) bool {
	var notFoundError *NotFoundError
	ok := errors.As(target, &notFoundError)
	return ok
}
