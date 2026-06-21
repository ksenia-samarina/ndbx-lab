package reactions

import (
	"context"
	"samarina/ndbx/internal/model"
)

type eventDomain interface {
	GetEventByID(ctx context.Context, eventId string) (model.Event, error)
	GetInternalUserID(ctx context.Context, sid model.Sid) (string, error)
}

type reactionsDomain interface {
	SetReaction(ctx context.Context, reaction model.Reaction) error
}
