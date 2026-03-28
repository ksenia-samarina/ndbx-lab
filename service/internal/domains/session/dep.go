package session

import (
	"context"
	"time"
)

type storage interface {
	CreateSession(ctx context.Context, sid *Sid, ttl time.Duration) error
	GetSession(ctx context.Context, sid *Sid) (bool, error)
	UpdateSession(ctx context.Context, sid *Sid, ttl time.Duration) error
}
