package validator

import "fmt"

type ErrInvalidFieldName struct {
	Field string
}

func (e *ErrInvalidFieldName) Error() string {
	return fmt.Sprintf("invalid %s field", e.Field)
}
