package storage

import (
	"errors"
	"fmt"
	"github.com/harshabangi/url-shortener/internal/storage/memory"
	"github.com/harshabangi/url-shortener/internal/storage/redis"
)

type Store interface {
	// SaveURL saves the mapping between a short URL key and the original URL.
	SaveURL(key string, originalURL string) error

	// GetOriginalURL retrieves the original URL associated with a short URL key.
	// It returns ErrNotFound if the key is not found.
	GetOriginalURL(key string) (string, error)
}

var ErrNotFound = errors.New("not found")

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
