package users

import (
	"context"
	"net/url"
	"samarina/ndbx/internal/model"
	"time"
)

type userDomain interface {
	RegisterUser(ctx context.Context, user model.User, ttl time.Duration) (model.Sid, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUsers(ctx context.Context, id, name string, limit, offset uint64) ([]model.User, error)
	GetUserByUserID(ctx context.Context, id string) (model.User, error)
	GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error)

	UpdateUserSession(ctx context.Context, sid model.Sid, ttl time.Duration) error
}

type validatorDomain interface {
	ValidateParams(query url.Values) (model.EventFilter, error)
}

type reactionsDomain interface {
	EnrichEventsWithReactions(ctx context.Context, events []model.Event, includeReactions bool) ([]model.Event, error)
}

type reviewsDomain interface {
	EnrichEventsWithReviews(ctx context.Context, events []model.Event, includeReviews bool) ([]model.Event, error)
}
