package events

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type domain interface {
	UpdateUserSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
	GetInternalUserID(ctx context.Context, sid model.Sid) (string, error)
	GetSession(ctx context.Context, sid model.Sid) (bool, error)

	GetEventsByUserID(ctx context.Context, id string) ([]model.Event, error)
	RegisterEvent(ctx context.Context, createdBy string, event model.Event) (string, error)
	UpdateEventsLocationCity(ctx context.Context, id string, city string) error
	GetEventByID(ctx context.Context, createdBy string) (model.Event, error)
	GetEvents(ctx context.Context, filters model.EventFilter) ([]model.Event, error)

	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}
