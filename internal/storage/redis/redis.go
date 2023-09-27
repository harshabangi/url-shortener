package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"sort"
)

type redisStore struct {
	c *redis.Client
}

func New(redisURL string) (*redisStore, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(options)

	// if we aren't able to ping redis, fail fast
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("could not connect to redis db: %w", err)
	}
	return &redisStore{c: client}, nil
}

func (m *redisStore) SaveURL(ctx context.Context, key, originalURL string) (string, error) {
	existingOriginalURL, err := m.c.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		_, err = m.c.Set(ctx, key, originalURL, 0).Result()
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

func (m *redisStore) GetOriginalURL(ctx context.Context, key string) (string, error) {
	value, err := m.c.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return "", shared.ErrNotFound
	case err != nil:
		return "", err
	default:
		return value, nil
	}
}

func (m *redisStore) RecordDomainFrequency(ctx context.Context, domainName string) error {
	_, err := m.c.ZIncrBy(ctx, "domain_frequencies", 1.0, domainName).Result()
	return err
}

func (m *redisStore) GetTopNDomainsByFrequency(ctx context.Context, n int) ([]shared.DomainFrequency, error) {
	result, err := m.c.ZRevRangeWithScores(ctx, "domain_frequencies", 0, int64(n-1)).Result()
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
