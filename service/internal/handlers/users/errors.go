package users

import (
	"errors"
	"fmt"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type ErrInvalidFieldName struct {
	Field string
}

func (e *ErrInvalidFieldName) Error() string {
	return fmt.Sprintf("invalid %s field", e.Field)
}

var ErrUserNotFound = errors.New("Not found")
