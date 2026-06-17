package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
	"time"
)

func (s *Storage) SetCounters(ctx context.Context, eventTitle string, counters *model.ReactionCounters, ttl time.Duration) error {
	key := s.buildReactionsKey(eventTitle)
	s.client.Del(ctx, key)
	err := s.client.HSet(ctx, key,
		"likes", counters.Likes,
		"dislikes", counters.Dislikes,
	).Err()
	if err != nil {
		return fmt.Errorf("redis hset error: %w", err)
	}
	return s.client.Expire(ctx, key, ttl).Err()
}
