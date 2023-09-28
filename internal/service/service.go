package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/harshabangi/url-shortener/docs"
	"github.com/harshabangi/url-shortener/internal/storage"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"github.com/harshabangi/url-shortener/internal/util"
	"github.com/harshabangi/url-shortener/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"strconv"
)

type Service struct {
	config  *Config
	storage storage.Store
}

type Config struct {
	ListenAddr        string `json:"listen_addr"`
	DataStorageEngine string `json:"data_storage_engine"`
	ShortURLDomain    string `json:"short_url_domain"`
	ShortURLLength    int    `json:"short_url_length"`
	RedisURL          string `json:"redis_url"`
}

func (c *Config) toStorageConfig() storage.Config {
	return storage.Config{
		DataStorageEngine: c.DataStorageEngine,
		RedisURL:          c.RedisURL,
	}
}

func NewConfig() *Config {
	return &Config{}
}

func NewService(config *Config) (*Service, error) {
	store, err := storage.New(config.toStorageConfig())
	if err != nil {
		return nil, err
	}

	return &Service{
		storage: store,
		config:  config,
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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowOrigins: []string{"*"},
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/v1/shorten", errorLoggingMiddleware(shorten))
	e.GET("/:short_code", errorLoggingMiddleware(expand))
	e.GET("/v1/metrics", errorLoggingMiddleware(metrics))

	e.Logger.Fatal(e.Start(s.config.ListenAddr))
}

// identify godoc
// @Summary Shorten the given long URL.
// @Description Shorten the given long URL.
// @Tags root
// @Param contact body pkg.ShortenRequest true "Shorten Request Body"
// @Accept json
// @Produce json
// @Success 200 {object} pkg.ShortenResponse
// @Failure 404
// @Failure 500
// @Router /v1/shorten [post]
// @Consumes application/json
func shorten(c echo.Context) error {
	s := c.Get("service").(*Service)

	var req pkg.ShortenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validationResult := req.Validate()
	if validationResult.Err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationResult.Err)
	}

	ctx := context.Background()

	shortCode, err := generateShortCode(ctx, s, req.URL, validationResult.Domain, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	response := &pkg.ShortenResponse{
		URL: fmt.Sprintf("%s/%s", s.config.ShortURLDomain, shortCode),
	}
	return c.JSON(http.StatusOK, response)
}

func generateShortCode(ctx context.Context, s *Service, originalLongURL, domainName string, collisionCounter int64) (string, error) {
	var (
		input   = []byte(originalLongURL)
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

	shortCode := hash[:s.config.ShortURLLength]

	existingLongURL, err := s.storage.SaveURL(ctx, shortCode, originalLongURL)
	switch {
	case err == nil:
		if err = s.storage.RecordDomainFrequency(ctx, domainName); err != nil {
			return "", err
		}
		return shortCode, nil

	case errors.Is(err, shared.ErrCollision):
		if originalLongURL == existingLongURL {
			return shortCode, nil
		}
		return generateShortCode(ctx, s, originalLongURL, domainName, collisionCounter+1)

	default:
		return "", err
	}
}

// @Summary Redirect to the original URL given a short code.
// @Description Redirect to the original URL associated with the provided short code.
// @Tags root
// @Accept json
// @Param short_code path string true "The short code to expand"
// @Success 301
// @Failure 404
// @Failure 500
// @Router /{short_code} [get]
func expand(c echo.Context) error {
	s := c.Get("service").(*Service)
	key := c.Param("short_code")
	originalLongURL, err := s.storage.GetOriginalURL(context.Background(), key)

	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Short code not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Redirect(http.StatusMovedPermanently, originalLongURL)
}

// @Summary Retrieve the top N domain names with the highest frequency of shortening.
// @Description Returns the top N domain names that have been shortened the most number of times.
// @Tags root
// @Accept json
// @Produce json
// @Param limit query int false "Number of top domains to retrieve (default: 3)"
// @Success 200 {object} pkg.DomainFreqListResponse
// @Failure 500
// @Router /v1/metrics [get]
func metrics(c echo.Context) error {
	s := c.Get("service").(*Service)
	limit := deriveLimit(c.QueryParam("limit"))
	domainFrequencies, err := s.storage.GetTopNDomainsByFrequency(context.Background(), limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	response := toDomainFreqListResponse(domainFrequencies)
	return c.JSON(http.StatusOK, &response)
}

func deriveLimit(limitParam string) int {
	limit := 3 // Default value
	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil {
			limit = l
		}
	}
	return limit
}

func toDomainFreqListResponse(frequencies []shared.DomainFrequency) pkg.DomainFreqListResponse {
	res := pkg.DomainFreqListResponse{}
	for _, v := range frequencies {
		res = append(res, pkg.DomainFreqResponse{DomainName: v.Domain, Frequency: v.Frequency})
	}
	return res
}

func errorLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		return err
	}
}
