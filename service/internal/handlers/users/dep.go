package users

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"
)

type domain interface {
	RegisterUser(ctx context.Context, user *types.User, ttl time.Duration) (*types.Sid, error)
	GetByUsername(ctx context.Context, username string) (*types.User, error)
	UpdateUserSession(ctx context.Context, sid *types.Sid, ttl time.Duration) error
}
