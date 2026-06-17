package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) GetCounters(ctx context.Context, eventTitle string) (*model.ReactionCounters, error) {
	key := s.buildReactionsKey(eventTitle)

	res, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis hgetall counters error: %w", err)
	}

	if len(res) == 0 {
		return nil, redis.Nil
	}

	var counters model.ReactionCounters

	if val, ok := res["likes"]; ok {
		_, err := fmt.Sscanf(val, "%d", &counters.Likes)
		if err != nil {
			return nil, err
		}
	}
	if val, ok := res["dislikes"]; ok {
		_, err := fmt.Sscanf(val, "%d", &counters.Dislikes)
		if err != nil {
			return nil, err
		}
	}

	return &counters, nil
}
