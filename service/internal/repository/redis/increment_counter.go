package redis

import (
	"context"
	"fmt"
	"time"
)

func (s *Storage) IncrementCounter(ctx context.Context, eventTitle string, field string, value int64, ttl time.Duration) error {
	key := s.buildKey(eventTitle)

	pipe := s.client.Pipeline()

	pipe.HSetNX(ctx, key, "likes", "0")
	pipe.HSetNX(ctx, key, "dislikes", "0")

	pipe.HIncrBy(ctx, key, field, value)

	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis atomic increment pipeline error: %w", err)
	}

	return nil
}
