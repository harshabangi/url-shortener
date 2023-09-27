package service

import (
	"github.com/harshabangi/url-shortener/internal/storage/shared"
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

func (ms *mockStorage) RecordDomainFrequency(domainName string) error {
	args := ms.Called(domainName)
	return args.Error(0)
}

func (ms *mockStorage) GetTopNDomainsByFrequency(n int) ([]shared.DomainFrequency, error) {
	args := ms.Called(n)
	return args.Get(0).([]shared.DomainFrequency), args.Error(1)
}
