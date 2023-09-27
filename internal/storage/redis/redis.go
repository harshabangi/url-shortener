package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"sort"
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
	existingOriginalURL, err := m.c.Get(key).Result()
	switch {
	case err == redis.Nil:
		_, err = m.c.Set(key, originalURL, 0).Result()
		if err != nil {
			return "", err
		}
		return "", nil
	case err != nil:
		return "", err
	default:
		return existingOriginalURL, shared.ErrCollision
	}
}

func (m *redisStore) GetOriginalURL(key string) (string, error) {
	value, err := m.c.Get(key).Result()
	switch {
	case err == redis.Nil:
		return "", shared.ErrNotFound
	case err != nil:
		return "", err
	default:
		return value, nil
	}
}

func (m *redisStore) RecordDomainFrequency(domainName string) error {
	_, err := m.c.ZIncrBy("domain_frequencies", 1.0, domainName).Result()
	return err
}

func (m *redisStore) GetTopNDomainsByFrequency(n int) ([]shared.DomainFrequency, error) {
	result, err := m.c.ZRevRangeWithScores("domain_frequencies", 0, int64(n-1)).Result()
	if err != nil {
		return nil, err
	}

	domainFrequencies := make([]shared.DomainFrequency, len(result))
	for i, v := range result {
		domainFrequencies[i] = shared.DomainFrequency{
			Domain:    v.Member.(string),
			Frequency: int64(v.Score),
		}
	}

	sort.Slice(domainFrequencies, func(i, j int) bool {
		return domainFrequencies[i].Frequency > domainFrequencies[j].Frequency
	})

	return domainFrequencies, nil
}
