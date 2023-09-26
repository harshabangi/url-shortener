package service

import (
	"github.com/harshabangi/url-shortener/internal/storage"
	"github.com/labstack/echo/v4"
	"os"
)

type Service struct {
	storage storage.Store
}

func NewService() (*Service, error) {
	store, err := storage.New("")
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

	e.Logger.Fatal(e.Start(os.Getenv("LISTEN_ADDR")))
}

func shorten(c echo.Context) error {
	return nil
}

func expand(c echo.Context) error {
	return nil
}

func metrics(c echo.Context) error {
	return nil
}
