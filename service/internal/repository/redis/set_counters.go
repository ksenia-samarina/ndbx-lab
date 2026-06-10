package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
	"time"
)

func (s *Storage) SetCounters(ctx context.Context, eventTitle string, counters *model.ReactionCounters, ttl time.Duration) error {
	key := s.buildKey(eventTitle)

	err := s.client.HSet(ctx, key,
		"likes", counters.Likes,
		"dislikes", counters.Dislikes,
	).Err()
	if err != nil {
		return fmt.Errorf("redis hset counters error: %w", err)
	}

	err = s.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis expire counters error: %w", err)
	}

	return nil
}
