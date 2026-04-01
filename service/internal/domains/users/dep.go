package users

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"
)

type sessionStorage interface {
	CreateUserSession(ctx context.Context, username string, sid *types.Sid, ttl time.Duration) error
	UpdateSession(ctx context.Context, sid *types.Sid, ttl time.Duration) error
}

type usersStorage interface {
	RegisterUser(ctx context.Context, user *types.User) error
	GetByUsername(ctx context.Context, username string) (*types.User, error)
}
