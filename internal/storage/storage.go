package storage

import (
	"fmt"
	"github.com/harshabangi/url-shortener/internal/storage/memory"
	"github.com/harshabangi/url-shortener/internal/storage/redis"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
)

type Store interface {
	// SaveURL saves the mapping between a short URL key and the original URL and returns the already existing original URL if it finds a match
	SaveURL(key, originalURL string) (string, error)

	// GetOriginalURL retrieves the original URL associated with a short URL key.
	// It returns ErrNotFound if the key is not found.
	GetOriginalURL(key string) (string, error)

	// RecordDomainFrequency stores the frequency of a domain name.
	// It associates the domainName with its frequency in the storage.
	RecordDomainFrequency(domainName string) error

	// GetTopNDomainsByFrequency returns the top n domains and their respective frequencies ordered by frequency descending.
	GetTopNDomainsByFrequency(n int) ([]shared.DomainFrequency, error)
}

func New(dataStorageEngine string) (Store, error) {
	switch dataStorageEngine {
	case "memory":
		return memory.New(), nil
	case "redis":
		return redis.New("", "")
	default:
		return nil, fmt.Errorf("storage engine '%s' is unsupported", dataStorageEngine)
	}
}
