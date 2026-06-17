package reviews

import (
	"context"
	"math"
	"samarina/ndbx/internal/model"
	"samarina/ndbx/internal/repository/cassandra"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (d *Domain) AddReview(ctx context.Context, eventID, createdBy, comment string, rating int) (string, error) {
	if rating < 1 || rating > 5 {
		return "", &ErrInvalidFieldName{"rating"}
	}
	if len(comment) == 0 || len(comment) > 300 {
		return "", &ErrInvalidFieldName{"comment"}
	}

	currentEvent, err := d.eventsStorage.GetEventByID(ctx, eventID)
	if err != nil {
		return "", ErrEventNotFound
	}

	existingReviews, _ := d.reviewStorage.GetReviewsByEventID(ctx, eventID)
	for _, r := range existingReviews {
		if r.CreatedBy == createdBy {
			return "", cassandra.ErrAlreadyExists
		}
	}

	review := model.Review{
		ID:        gocql.TimeUUID().String(),
		EventID:   eventID,
		CreatedBy: createdBy,
		Comment:   comment,
		Rating:    rating,
		CreatedAt: time.Now().UTC(),
	}

	reviewID, err := d.reviewStorage.CreateReview(ctx, review)
	if err != nil {
		return "", err
	}

	allEvents, err := d.eventsStorage.GetEvents(ctx, model.EventFilter{Title: currentEvent.Title})
	if err != nil || len(allEvents) == 0 {
		allEvents = []model.Event{currentEvent}
	}

	var totalRating float64
	var totalCount int

	for _, ev := range allEvents {
		idSearch := ev.ID.Hex()
		reviewsList, cassErr := d.reviewStorage.GetReviewsByEventID(ctx, idSearch)
		if cassErr == nil {
			for _, r := range reviewsList {
				totalRating += float64(r.Rating)
				totalCount++
			}
		}
	}

	avgRating := 0.0
	if totalCount > 0 {
		avgRating = totalRating / float64(totalCount)
		avgRating = math.Round(avgRating*10) / 10
	}

	newSummary := &model.ReviewsSummary{Count: totalCount, Rating: avgRating}
	_ = d.reviewsCache.SetReviewsSummary(ctx, currentEvent.Title, newSummary, d.reviewTTL)

	return reviewID, nil
}
