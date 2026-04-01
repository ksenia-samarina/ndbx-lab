package session

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"
)

type domain interface {
	UpdateSession(ctx context.Context, sid *types.Sid, ttl time.Duration) error
	CreateSession(ctx context.Context, ttl time.Duration) (*types.Sid, error)
	GetSession(ctx context.Context, sid *types.Sid) (bool, error)
}
