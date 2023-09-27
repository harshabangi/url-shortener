package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/harshabangi/url-shortener/internal/storage"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"github.com/harshabangi/url-shortener/internal/util"
	"github.com/harshabangi/url-shortener/pkg"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service struct {
	storage storage.Store
}

func NewService() (*Service, error) {
	store, err := storage.New("memory")
	if err != nil {
		return nil, err
	}

	return &Service{
		storage: store,
	}, nil
}

func (s *Service) Run() {
	e := echo.New()

	// Register app (*App) to be injected into all HTTP handlers.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("service", s)
			return next(c)
		}
	})

	//e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/v1/shorten", shorten)
	e.GET("/v1/expand/:short_code", expand)
	e.GET("/v1/metrics", metrics)

	e.Logger.Fatal(e.Start(":8082"))
}

func shorten(c echo.Context) error {
	s := c.Get("service").(*Service)

	var req pkg.ShortenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validationResult := req.Validate()
	if validationResult.Err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	shortURL, err := createShortURLHash(s, req.URL, 0)
	if err == nil {
		return echo.NewHTTPError(http.StatusOK, &pkg.ShortenResponse{URL: shortURL})
	}
	return echo.NewHTTPError(http.StatusInternalServerError)
}

func createShortURLHash(s *Service, originalURL string, collisionCounter int64) (string, error) {
	var (
		input   = []byte(originalURL)
		counter = []byte(fmt.Sprintf("%d", collisionCounter))
	)
	input = append(input, counter...)

	hasher := md5.New()
	hasher.Write(input)
	md5Hash := hex.EncodeToString(hasher.Sum(nil))

	hash, err := util.Md5ToBase62(md5Hash)
	if err != nil {
		return "", fmt.Errorf("error converting from md5 to base62: %s: %w", md5Hash, err)
	}

	shortHash := hash[:7]

	err = s.storage.SaveURL(shortHash, originalURL)
	switch {
	case err == nil:
		return shortHash, nil
	case errors.Is(err, shared.ErrCollision):
		return createShortURLHash(s, originalURL, collisionCounter+1)
	default:
		return "", err
	}
}

func expand(c echo.Context) error {
	return nil
}

func metrics(c echo.Context) error {
	return nil
}
