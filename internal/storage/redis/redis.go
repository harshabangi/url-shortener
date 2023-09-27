package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type redisStore struct {
	c *redis.Client
}

func New(hostAddr, password string) (*redisStore, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     hostAddr,
		Password: password,
	})

	// if we aren't able to ping redis, fail fast
	if _, err := c.Ping().Result(); err != nil {
		return nil, fmt.Errorf("could not connect to redis db: %w", err)
	}
	return &redisStore{c: c}, nil
}

func (m *redisStore) SaveURL(key, originalURL string) error {
	//TODO implement me
	panic("implement me")
}

func (m *redisStore) GetOriginalURL(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *redisStore) GetShortURLKey(originalURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}
