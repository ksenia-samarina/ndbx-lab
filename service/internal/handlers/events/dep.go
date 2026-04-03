package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
	"time"
)

type domain interface {
	GetSession(ctx context.Context, sid *types.Sid) (bool, error)
	GetInternalUserID(ctx context.Context, sid *types.Sid) (string, error)
	GetEvent(ctx context.Context, createdBy string) (*types.Event, error)
	CreateEvent(ctx context.Context, createdBy string, event *types.Event) (string, error)
	ListEvents(ctx context.Context, title string, offset int64, limit int64) ([]types.Event, error)
	UpdateUserSession(ctx context.Context, sid *types.Sid, ttl time.Duration) error
}
