package reviews

import (
	"errors"
	"fmt"
)

type ErrInvalidFieldName struct {
	Field string
}

func (e *ErrInvalidFieldName) Error() string {
	return fmt.Sprintf("invalid %s field", e.Field)
}

var ErrEventNotFound = errors.New("event not found")
var ErrNoReviewAccess = errors.New("no review access")
