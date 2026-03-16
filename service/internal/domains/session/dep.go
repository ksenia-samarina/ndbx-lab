package session

import (
	"context"
	"time"
)

type Storage interface {
	CreateUserSession(ctx context.Context, id Id, expiration time.Duration) error
	UpdateUserSession(ctx context.Context, id Id, expiration time.Duration) error
	CheckExistenceSession(ctx context.Context, id Id) (bool, error)
}
