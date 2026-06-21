package reactions

import (
	"context"
	"samarina/ndbx/internal/model"
)

type eventDomain interface {
	GetEventByID(ctx context.Context, eventId string) (model.Event, error)
	GetInternalUserID(ctx context.Context, sid model.Sid) (string, error)
}

type reviewDomain interface {
	AddReview(ctx context.Context, eventID, createdBy, comment string, rating int) (string, error)
	GetReviews(ctx context.Context, eventID string, limit, offset int) ([]model.Review, int, error)
	PatchReview(ctx context.Context, eventID, reviewID, userID string, newRating *int, newComment *string) error
}
