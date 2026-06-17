package reviews

import (
	"context"
	"math"
	"samarina/ndbx/internal/model"
	"time"
)

func (d *Domain) EnrichEventsWithReviews(ctx context.Context, events []model.Event, includeReviews bool) ([]model.Event, error) {
	if !includeReviews {
		return events, nil
	}

	for i := range events {
		defaultSummary := &model.ReviewsSummary{Count: 0, Rating: 0.0}
		title := events[i].Title

		summary, err := d.reviewsCache.GetReviewsSummary(ctx, title)
		if err == nil && summary != nil {
			events[i].Reviews = summary
			continue
		}

		allEventsWithSameTitle, dbErr := d.eventsStorage.GetEvents(ctx, model.EventFilter{
			Title: title,
		})
		if dbErr != nil || len(allEventsWithSameTitle) == 0 {
			allEventsWithSameTitle = []model.Event{events[i]}
		}

		var totalRating float64
		var totalCount int

		for _, ev := range allEventsWithSameTitle {
			idSearch := ev.ID.Hex()

			reviewsList, cassErr := d.reviewStorage.GetReviewsByEventID(ctx, idSearch)
			if cassErr == nil {
				for _, r := range reviewsList {
					totalRating += float64(r.Rating)
					totalCount++
				}
			}
		}

		if totalCount == 0 {
			events[i].Reviews = defaultSummary
			_ = d.reviewsCache.SetReviewsSummary(ctx, title, defaultSummary, 1*time.Minute)
			continue
		}

		avgRating := totalRating / float64(totalCount)
		avgRating = math.Round(avgRating*10) / 10

		calculatedSummary := &model.ReviewsSummary{
			Count:  totalCount,
			Rating: avgRating,
		}

		_ = d.reviewsCache.SetReviewsSummary(ctx, title, calculatedSummary, d.reviewTTL)
		events[i].Reviews = calculatedSummary
	}

	return events, nil
}
