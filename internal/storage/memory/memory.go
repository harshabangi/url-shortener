package memory

import (
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"sync"
)

func New() *memoryStore {
	return &memoryStore{
		urls:  make(map[string]string),
		mutex: sync.RWMutex{},
	}
}

type memoryStore struct {
	urls  map[string]string // Maps short codes to original URLs
	mutex sync.RWMutex      // Mutex for thread-safe access
}

func (m *memoryStore) SaveURL(key, originalURL string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.urls[key]; ok {
		return shared.ErrCollision
	}
	m.urls[key] = originalURL
	return nil
}

func (m *memoryStore) GetOriginalURL(key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if value, ok := m.urls[key]; ok {
		return value, nil
	}
	return "", shared.ErrNotFound
}
