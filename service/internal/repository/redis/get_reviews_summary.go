package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) GetReviewsSummary(ctx context.Context, eventTitle string) (*model.ReviewsSummary, error) {
	key := s.buildReviewsKey(eventTitle)

	res, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis hgetall reviews error: %w", err)
	}

	if len(res) == 0 {
		return nil, redis.Nil
	}

	var summary model.ReviewsSummary

	if val, ok := res["count"]; ok {
		summary.Count, _ = strconv.Atoi(val)
	}
	if val, ok := res["rating"]; ok {
		summary.Rating, _ = strconv.ParseFloat(val, 64)
	}

	return &summary, nil
}
