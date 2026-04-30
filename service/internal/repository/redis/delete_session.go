package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/domains/types"
)

func (s *Storage) DeleteSession(ctx context.Context, sid *types.Sid) error {
	err := s.client.Del(ctx, sid.SidString).Err()
	if err != nil {
		return fmt.Errorf("delete session error: %w", err)
	}
	return nil
}
