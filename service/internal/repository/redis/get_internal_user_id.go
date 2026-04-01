package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/domains/types"
)

func (s *Storage) GetInternalUserID(ctx context.Context, sid *types.Sid) (string, error) {
	userID, err := s.client.HGet(ctx, sid.SidString, "user_id").Result()
	if err != nil {
		return "", fmt.Errorf("error get user_id from user session: %w", err)
	}
	return userID, nil
}
