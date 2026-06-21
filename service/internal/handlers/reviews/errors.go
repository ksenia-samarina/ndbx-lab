package reactions

import (
	"errors"
	"fmt"
)

var ErrEventNotFound = errors.New("Event not found")
var ErrAlreadyExists = errors.New("Already exists")

type ErrInvalidFieldName struct {
	Field string
}

func (e *ErrInvalidFieldName) Error() string {
	return fmt.Sprintf("invalid %s field", e.Field)
}
