package reactions

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type reactionsStorage interface {
	SetReaction(ctx context.Context, reaction model.Reaction) error
	GetReactionCounters(ctx context.Context, eventID string) (*model.ReactionCounters, error)
}

type reactionsCache interface {
	GetCounters(ctx context.Context, eventTitle string) (*model.ReactionCounters, error)
	SetCounters(ctx context.Context, eventTitle string, counters *model.ReactionCounters, ttl time.Duration) error
}

type eventsStorage interface {
	GetEventByID(ctx context.Context, id string) (model.Event, error)
	GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error)
}
