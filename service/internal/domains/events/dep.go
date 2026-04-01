package events

import (
	"context"
	"samarina/ndbx/internal/domains/types"
)

type sessionStorage interface {
	GetInternalUserID(ctx context.Context, sid *types.Sid) (string, error)
}

type eventsStorage interface {
	GetEvent(ctx context.Context, createdBy string) (*types.Event, error)
	CreateEvent(ctx context.Context, createdBy string, event *types.Event) (string, error)
	ListEvents(ctx context.Context, title string, offset int64, limit int64) ([]types.Event, error)
}
