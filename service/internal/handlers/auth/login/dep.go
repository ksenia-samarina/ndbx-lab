package login

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"
)

type domain interface {
	GetByUsername(ctx context.Context, username string) (*types.User, error)
	GetInternalUserID(ctx context.Context, username string) (string, error)
	GetSession(ctx context.Context, sid *types.Sid) (bool, error)
	UpdateUserSession(ctx context.Context, userID string, sid *types.Sid, ttl time.Duration) error
	CreateUserSession(ctx context.Context, userID string, ttl time.Duration) (*types.Sid, error)
}
