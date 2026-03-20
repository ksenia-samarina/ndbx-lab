package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/domains/session"
)

func (s *SessionStorage) GetSession(ctx context.Context, sid *session.Sid) (bool, error) {
	n, err := s.client.Exists(ctx, sid.SidString).Result()
	if err != nil {
		return false, fmt.Errorf("error check session existence: %w", err)
	}
	return n > 0, nil
}
