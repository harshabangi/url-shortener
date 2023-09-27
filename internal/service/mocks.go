package service

import (
	"github.com/stretchr/testify/mock"
	"io"
)

type mockStorage struct {
	io.Closer
	mock.Mock
}

func (ms *mockStorage) SaveURL(key, originalURL string) error {
	args := ms.Called(key, originalURL)
	return args.Error(0)
}

func (ms *mockStorage) GetOriginalURL(key string) (string, error) {
	args := ms.Called(key)
	return args.Get(0).(string), args.Error(1)
}
