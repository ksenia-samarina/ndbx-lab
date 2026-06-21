package reviews

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"
)

type reviewStorage interface {
	CreateReview(ctx context.Context, r model.Review) (string, error)
	GetReviewsByEventID(ctx context.Context, eventID string) ([]model.Review, error)
	GetReviewByID(ctx context.Context, eventID, reviewID string) (*model.Review, error)
	UpdateReview(ctx context.Context, r model.Review) error
}

type eventsStorage interface {
	GetEventByID(ctx context.Context, id string) (model.Event, error)
	GetEvents(ctx context.Context, filter model.EventFilter) ([]model.Event, error)
}

type reviewsCache interface {
	GetReviewsSummary(ctx context.Context, eventID string) (*model.ReviewsSummary, error)
	SetReviewsSummary(ctx context.Context, eventID string, summary *model.ReviewsSummary, ttl time.Duration) error
}
