package redis

import (
	"context"
	"samarina/ndbx/internal/model"
)

func (s *Storage) GetInternalUserID(ctx context.Context, sid model.Sid) (string, error) {
	return s.client.HGet(ctx, sid.SidString, "user_id").Result()
}
