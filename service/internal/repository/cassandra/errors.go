package cassandra

import "errors"

var ErrAlreadyExists = errors.New("Already exists")
var ErrReviewNotFound = errors.New("Review not found")
