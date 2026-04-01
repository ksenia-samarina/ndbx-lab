package login

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sessionStorage interface {
	CreateUserSession(ctx context.Context, userID string, sid *types.Sid, ttl time.Duration) error
	GetSession(ctx context.Context, sid *types.Sid) (bool, error)
	UpdateUserSession(ctx context.Context, userID string, sid *types.Sid, ttl time.Duration) error
}

type loginStorage interface {
	GetByUsername(ctx context.Context, username string) (*types.User, error)
	GetInternalUserID(ctx context.Context, username string) (primitive.ObjectID, error)
}
