package login

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authStorage interface {
	CreateUserSession(ctx context.Context, userID string, sid model.Sid, ttl time.Duration) error
	GetSession(ctx context.Context, sid model.Sid) (bool, error)
	UpdateUserSession(ctx context.Context, userID string, sid model.Sid, ttl time.Duration) error
}

type loginStorage interface {
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetInternalUserID(ctx context.Context, username string) (primitive.ObjectID, error)
}
