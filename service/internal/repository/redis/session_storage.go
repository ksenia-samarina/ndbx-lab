package redis

import (
	"context"
	"fmt"
	"samarina/ndbx/internal/domains/session"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStorage struct {
	client *redis.Client
}

func (s *SessionStorage) CreateUserSession(ctx context.Context, id session.Id, expiration time.Duration) error {
	now := time.Now().UTC().Format(time.RFC3339)
	sid := "sid:" + id.HexString

	_, err := s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, sid, "created_at", now, "updated_at", now)
		pipe.Expire(ctx, sid, expiration)
		return nil
	})
	return err
}

func (s *SessionStorage) UpdateUserSession(ctx context.Context, id session.Id, expiration time.Duration) error {
	now := time.Now().UTC().Format(time.RFC3339)
	sid := "sid:" + id.HexString

	_, err := s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, sid, "updated_at", now)
		pipe.Expire(ctx, sid, expiration)
		return nil
	})
	return err
}

func (s *SessionStorage) CheckExistenceSession(ctx context.Context, id session.Id) (bool, error) {
	key := "sid:" + id.HexString
	n, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error check session existence: %w", err)
	}
	return n > 0, nil
}

func NewSessionStorage(client *redis.Client) *SessionStorage {
	return &SessionStorage{client: client}
}
