package login

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type domain interface {
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetInternalUserID(ctx context.Context, username string) (string, error)
	GetSession(ctx context.Context, sid model.Sid) (bool, error)
	UpdateUserSession(ctx context.Context, userID string, sid model.Sid, ttl time.Duration) error
	CreateUserSession(ctx context.Context, userID string, ttl time.Duration) (model.Sid, error)
}
