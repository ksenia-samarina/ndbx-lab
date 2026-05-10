package users

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sessionStorage interface {
	CreateUserSession(ctx context.Context, username string, sid model.Sid, ttl time.Duration) error
	UpdateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
}

type usersStorage interface {
	RegisterUser(ctx context.Context, user model.User) error
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUserByUserID(ctx context.Context, id string) (model.User, error)
	GetUsers(ctx context.Context, id, name string, limit, offset uint64) ([]model.User, error)
	GetInternalUserID(ctx context.Context, username string) (primitive.ObjectID, error)
}

type eventStorage interface {
	GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error)
}
