package session

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"
)

type storage interface {
	CreateSession(ctx context.Context, sid *types.Sid, ttl time.Duration) error
	GetSession(ctx context.Context, sid *types.Sid) (bool, error)
	UpdateSession(ctx context.Context, sid *types.Sid, ttl time.Duration) error
}
