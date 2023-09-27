package service

import (
	"context"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"github.com/stretchr/testify/mock"
	"io"
)

type mockStorage struct {
	io.Closer
	mock.Mock
}

func (ms *mockStorage) SaveURL(ctx context.Context, key, originalURL string) (string, error) {
	args := ms.Called(ctx, key, originalURL)
	return args.Get(0).(string), args.Error(1)
}

func (ms *mockStorage) GetOriginalURL(ctx context.Context, key string) (string, error) {
	args := ms.Called(ctx, key)
	return args.Get(0).(string), args.Error(1)
}

func (ms *mockStorage) RecordDomainFrequency(ctx context.Context, domainName string) error {
	args := ms.Called(ctx, domainName)
	return args.Error(0)
}

func (ms *mockStorage) GetTopNDomainsByFrequency(ctx context.Context, n int) ([]shared.DomainFrequency, error) {
	args := ms.Called(ctx, n)
	return args.Get(0).([]shared.DomainFrequency), args.Error(1)
}
