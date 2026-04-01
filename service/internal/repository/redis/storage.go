package redis

import "github.com/redis/go-redis/v9"

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
}
