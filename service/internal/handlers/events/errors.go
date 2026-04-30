package events

import (
	"errors"
	"fmt"
)

var ErrEventAlreadyExists = errors.New("event already exists")
var ErrEventNotExist = errors.New("Not found")

type ErrInvalidFieldName struct {
	Field string
}

func (e *ErrInvalidFieldName) Error() string {
	return fmt.Sprintf("invalid %s field", e.Field)
}
