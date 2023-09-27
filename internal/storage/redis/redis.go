package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
)

type redisStore struct {
	c *redis.Client
}

func New(hostAddr, password string) (*redisStore, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     hostAddr,
		Password: password,
		DB:       0,
	})

	// if we aren't able to ping redis, fail fast
	if _, err := c.Ping().Result(); err != nil {
		return nil, fmt.Errorf("could not connect to redis db: %w", err)
	}
	return &redisStore{c: c}, nil
}

func (m *redisStore) SaveURL(key, originalURL string) (string, error) {
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

func (m *redisStore) RecordDomainFrequency(domainName string) error {
	//TODO implement me
	panic("implement me")
}

func (m *redisStore) GetTopNDomainsByFrequency(n int) ([]shared.DomainFrequency, error) {
	//TODO implement me
	panic("implement me")
}
