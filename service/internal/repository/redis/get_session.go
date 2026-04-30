package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/domains/types"
)

func (s *Storage) GetSession(ctx context.Context, sid *types.Sid) (bool, error) {
	n, err := s.client.Exists(ctx, sid.SidString).Result()
	if err != nil {
		return false, fmt.Errorf("error check session existence: %w", err)
	}
	return n > 0, nil
}
