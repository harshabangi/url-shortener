package storage

import (
	"context"
	"fmt"
	"github.com/harshabangi/url-shortener/internal/storage/memory"
	"github.com/harshabangi/url-shortener/internal/storage/redis"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
)

type Store interface {
	// SaveURL saves the mapping between a short URL key and the original URL and returns the already existing original URL if it finds a match
	SaveURL(ctx context.Context, key, originalURL string) (string, error)

	// GetOriginalURL retrieves the original URL associated with a short URL key.
	// It returns ErrNotFound if the key is not found.
	GetOriginalURL(ctx context.Context, key string) (string, error)

	// RecordDomainFrequency stores the frequency of a domain name.
	// It associates the domainName with its frequency in the storage.
	RecordDomainFrequency(ctx context.Context, domainName string) error

	// GetTopNDomainsByFrequency returns the top n domains and their respective frequencies ordered by frequency descending.
	GetTopNDomainsByFrequency(ctx context.Context, n int) ([]shared.DomainFrequency, error)
}

type Config struct {
	DataStorageEngine string
	RedisURL          string
}

func New(config Config) (Store, error) {
	switch config.DataStorageEngine {
	case "memory":
		return memory.New(), nil
	case "redis":
		return redis.New(config.RedisURL)
	default:
		return nil, fmt.Errorf("storage engine '%s' is unsupported", config.DataStorageEngine)
	}
}
