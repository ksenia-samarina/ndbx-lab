package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/model"
)

func (s *Storage) DeleteSession(ctx context.Context, sid model.Sid) error {
	err := s.client.Del(ctx, sid.SidString).Err()
	if err != nil {
		return fmt.Errorf("delete auth error: %w", err)
	}
	return nil
}
