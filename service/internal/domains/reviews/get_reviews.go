package reviews

import (
	"context"
	"samarina/ndbx/internal/model"
)

import "sort"

func (d *Domain) GetReviews(ctx context.Context, eventID string, limit, offset int) ([]model.Review, int, error) {
	if limit < 0 {
		return nil, 0, &ErrInvalidFieldName{Field: `"limit"`}
	}
	if offset < 0 {
		return nil, 0, &ErrInvalidFieldName{Field: `"offset"`}
	}

	allReviews, err := d.reviewStorage.GetReviewsByEventID(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(allReviews, func(i, j int) bool {
		return allReviews[i].CreatedAt.After(allReviews[j].CreatedAt)
	})

	totalCount := len(allReviews)
	if offset > totalCount {
		return []model.Review{}, 0, nil
	}

	end := offset + limit
	if end > totalCount || limit == 0 {
		end = totalCount
	}
	return allReviews[offset:end], totalCount, nil
}
