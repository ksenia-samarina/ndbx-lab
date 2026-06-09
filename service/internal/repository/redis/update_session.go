package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) UpdateSession(ctx context.Context, sid model.Sid, ttl time.Duration) error {
	t := time.Now().UTC().Format(time.RFC3339)

	err := s.client.HSetEXWithArgs(ctx, sid.SidString,
		&redis.HSetEXOptions{
			ExpirationType: redis.HSetEXExpirationEX,
			ExpirationVal:  int64(ttl.Seconds()),
		},
		"updated_at", t,
	).Err()
	if err != nil {
		return fmt.Errorf("create session error: %w", err)
	}
	return nil
}
