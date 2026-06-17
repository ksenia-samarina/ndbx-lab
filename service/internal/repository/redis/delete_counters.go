package redis

import (
	"context"
	"fmt"
)

func (s *Storage) DeleteCounters(ctx context.Context, eventTitle string) error {
	key := s.buildReactionsKey(eventTitle)

	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete counters error: %w", err)
	}

	return nil
}
