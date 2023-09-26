package memory

import "github.com/harshabangi/url-shortener/internal/storage"

func New() storage.Store {
	return &memoryStore{}
}

type memoryStore struct {
	counter int64
}

func (m *memoryStore) SaveURL(key string, originalURL string) error {

}

func (m *memoryStore) GetOriginalURL(key string) (string, error) {

}
