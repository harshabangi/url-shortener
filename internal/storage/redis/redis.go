package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/harshabangi/url-shortener/internal/storage"
)

type redisStore struct {
	c *redis.Client
}

func New(hostAddr, password string) (storage.Store, error) {
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

func (m *redisStore) SaveURL(key string, originalURL string) error {

}

func (m *redisStore) GetOriginalURL(key string) (string, error) {

}
