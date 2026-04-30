package events

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type sessionStorage interface {
	UpdateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
	GetInternalUserID(ctx context.Context, sid model.Sid) (string, error)
}

type eventsStorage interface {
	GetEventByID(ctx context.Context, id string) (model.Event, error)
	GetEventsByUserID(ctx context.Context, createdBy string) ([]model.Event, error)
	GetEvents(ctx context.Context, filters model.EventFilter) ([]model.Event, error)
	RegisterEvent(ctx context.Context, createdBy string, event model.Event) (string, error)
	UpdateEventsLocationCity(ctx context.Context, id string, city string) error
}
