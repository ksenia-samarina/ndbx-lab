package session

import (
	"context"
	"samarina/ndbx/internal/domains/session"
	"time"
)

type domain interface {
	UpdateSession(ctx context.Context, sid *session.Sid, ttl time.Duration) error
	CreateSession(ctx context.Context, ttl time.Duration) (*session.Sid, error)
	GetSession(ctx context.Context, sid *session.Sid) (bool, error)
}
