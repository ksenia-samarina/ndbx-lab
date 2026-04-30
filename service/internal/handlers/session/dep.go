package session

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type domain interface {
	UpdateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
	CreateSession(ctx context.Context, ttl time.Duration) (model.Sid, error)
	GetSession(ctx context.Context, sid model.Sid) (bool, error)
}
