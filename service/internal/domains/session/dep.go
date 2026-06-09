package session

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type storage interface {
	CreateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
	GetSession(ctx context.Context, sid model.Sid) (bool, error)
	UpdateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
}
