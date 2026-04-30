package redis

import (
	"context"
	"errors"
	"fmt"
	"samarina/ndbx/internal/domains/types"

	redisdb "github.com/redis/go-redis/v9"
)

func (s *Storage) GetInternalUserID(ctx context.Context, sid *types.Sid) (string, error) {
	userID, err := s.client.HGet(ctx, sid.SidString, "user_id").Result()
	if errors.Is(err, redisdb.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("error get user_id from user session: %w", err)
	}
	return userID, nil
}
