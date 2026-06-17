package reviews

import (
	"context"
	"math"
	"samarina/ndbx/internal/model"
	"time"
)

func (d *Domain) PatchReview(ctx context.Context, eventID, reviewID, userID string, newRating *int, newComment *string) error {
	currentEvent, err := d.eventsStorage.GetEventByID(ctx, eventID)
	if err != nil {
		return ErrEventNotFound
	}

	review, err := d.reviewStorage.GetReviewByID(ctx, eventID, reviewID)
	if err != nil {
		return ErrEventNotFound
	}

	if review.CreatedBy != userID {
		return ErrNoReviewAccess
	}

	if newRating != nil {
		if *newRating < 1 || *newRating > 5 {
			return &ErrInvalidFieldName{Field: `"rating"`}
		}
		review.Rating = *newRating
	}

	if newComment != nil {
		if len(*newComment) == 0 || len(*newComment) > 300 {
			return &ErrInvalidFieldName{Field: `"comment"`}
		}
		review.Comment = *newComment
	}

	review.UpdatedAt = time.Now().UTC()

	err = d.reviewStorage.UpdateReview(ctx, *review)
	if err != nil {
		return err
	}

	return d.syncReviewsCache(ctx, currentEvent.Title)
}

func (d *Domain) syncReviewsCache(ctx context.Context, eventTitle string) error {
	allEvents, err := d.eventsStorage.GetEvents(ctx, model.EventFilter{
		Title: eventTitle,
	})
	if err != nil || len(allEvents) == 0 {
		return nil
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

	newSummary := &model.ReviewsSummary{
		Count:  totalCount,
		Rating: avgRating,
	}

	return d.reviewsCache.SetReviewsSummary(ctx, eventTitle, newSummary, d.reviewTTL)
}
