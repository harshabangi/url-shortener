package service

import (
	"context"
	"encoding/json"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"github.com/harshabangi/url-shortener/pkg"
	"github.com/labstack/echo/v4"
	asserts "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testService(ms *mockStorage) *Service {
	return &Service{
		storage: ms,
		config: &Config{
			DataStorageEngine: "memory",
			ShortURLDomain:    "http://localhost:8080",
			ShortURLLength:    7,
		},
	}
}

func Test_Shorten(t *testing.T) {

	t.Run("shorten the original long URL", func(t *testing.T) {
		assert := asserts.New(t)

		ms := &mockStorage{}
		mockService := testService(ms)

		rqBody := `{"url":"https://www.google.com"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/shorten", strings.NewReader(rqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.Set("service", mockService)
		ctx := context.Background()

		ms.On("SaveURL", ctx, "g2GJ99W", "https://www.google.com").Return("", nil)
		ms.On("RecordDomainFrequency", ctx, "www.google.com").Return(nil)

		want := &pkg.ShortenResponse{
			URL: "http://localhost:8080/g2GJ99W",
		}

		wantBytes, _ := json.Marshal(want)
		err := shorten(c)
		assert.Nil(err)
		assert.Equal(string(wantBytes), strings.Trim(rec.Body.String(), "\n"))

		ms.AssertExpectations(t)
	})

	t.Run("original long URL already exists", func(t *testing.T) {
		assert := asserts.New(t)

		ms := &mockStorage{}
		mockService := testService(ms)

		rqBody := `{"url":"https://www.google.com"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/shorten", strings.NewReader(rqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.Set("service", mockService)
		ctx := context.Background()

		ms.On("SaveURL", ctx, "g2GJ99W", "https://www.google.com").Return("https://www.google.com", shared.ErrCollision)

		want := &pkg.ShortenResponse{
			URL: "http://localhost:8080/g2GJ99W",
		}

		wantBytes, _ := json.Marshal(want)
		err := shorten(c)
		assert.Nil(err)
		assert.Equal(string(wantBytes), strings.Trim(rec.Body.String(), "\n"))

		ms.AssertExpectations(t)
	})

	t.Run("collision with another long URL", func(t *testing.T) {
		assert := asserts.New(t)

		ms := &mockStorage{}
		mockService := testService(ms)

		rqBody := `{"url":"https://www.google.com"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/shorten", strings.NewReader(rqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.Set("service", mockService)
		ctx := context.Background()

		ms.On("SaveURL", ctx, "g2GJ99W", "https://www.google.com").Return("https://www.google.com/1", shared.ErrCollision)
		ms.On("SaveURL", ctx, "baYj37f", "https://www.google.com").Return("", nil)
		ms.On("RecordDomainFrequency", ctx, "www.google.com").Return(nil)

		want := &pkg.ShortenResponse{
			URL: "http://localhost:8080/baYj37f",
		}

		wantBytes, _ := json.Marshal(want)
		err := shorten(c)
		assert.Nil(err)
		assert.Equal(string(wantBytes), strings.Trim(rec.Body.String(), "\n"))

		ms.AssertExpectations(t)
	})

}

func TestExpandHandler(t *testing.T) {

	t.Run("Short code found", func(t *testing.T) {
		assert := asserts.New(t)
		ms := &mockStorage{}
		mockService := testService(ms)

		req := httptest.NewRequest(http.MethodGet, "/:short_code", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("service", mockService)
		c.SetParamNames("short_code")
		c.SetParamValues("6sWZzzm")
		ctx := context.Background()

		ms.On("GetOriginalURL", ctx, "6sWZzzm").Return("https://www.google.com", nil)

		err := expand(c)
		assert.Nil(err)
		assert.Equal(http.StatusFound, rec.Code)
		assert.Equal("https://www.google.com", rec.Header().Get("Location"))

		ms.AssertExpectations(t)
	})

	t.Run("Short code not found", func(t *testing.T) {
		assert := asserts.New(t)
		ms := &mockStorage{}
		mockService := testService(ms)

		req := httptest.NewRequest(http.MethodGet, "/:short_code", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("service", mockService)
		c.SetParamNames("short_code")
		c.SetParamValues("6sWZzzm")
		ctx := context.Background()

		ms.On("GetOriginalURL", ctx, "6sWZzzm").Return("", shared.ErrNotFound)

		err := expand(c)
		assert.NotNil(err)

		ms.AssertExpectations(t)
	})
}

func TestDeriveLimit(t *testing.T) {
	assert := asserts.New(t)

	tcc := []struct {
		limitParam string
		expected   int
	}{
		{"", 3},
		{"5", 5},
		{"abc", 3},
		{"0", 0},
		{"100", 100},
	}

	for _, test := range tcc {
		t.Run(test.limitParam, func(t *testing.T) {
			limit := deriveLimit(test.limitParam)
			assert.Equal(test.expected, limit)
		})
	}
}
