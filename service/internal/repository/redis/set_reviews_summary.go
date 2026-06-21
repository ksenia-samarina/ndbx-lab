package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
	"strconv"
	"time"
)

func (s *Storage) SetReviewsSummary(ctx context.Context, eventTitle string, summary *model.ReviewsSummary, ttl time.Duration) error {
	key := s.buildReviewsKey(eventTitle)

	pipe := s.client.Pipeline()
	pipe.HSet(ctx, key, "count", strconv.Itoa(summary.Count))
	pipe.HSet(ctx, key, "rating", strconv.FormatFloat(summary.Rating, 'f', 1, 64))
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline set reviews hash error: %w", err)
	}
	return nil
}
