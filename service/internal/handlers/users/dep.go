package users

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type domain interface {
	RegisterUser(ctx context.Context, user model.User, ttl time.Duration) (model.Sid, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUsers(ctx context.Context, id, name string, limit, offset uint64) ([]model.User, error)
	GetUserByUserID(ctx context.Context, id string) (model.User, error)
	GetUserEventsByUserID(ctx context.Context, userID string) ([]model.Event, error)

	UpdateUserSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
}
