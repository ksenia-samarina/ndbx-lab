package redis

import "github.com/redis/go-redis/v9"

type SessionStorage struct {
	client *redis.Client
}

func NewSessionStorage(client *redis.Client) *SessionStorage {
	return &SessionStorage{client: client}
}
