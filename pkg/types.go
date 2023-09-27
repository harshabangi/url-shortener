package pkg

import (
	"errors"
	"fmt"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ValidationResult struct {
	Domain string
	Err    error
}

func (r ShortenRequest) Validate() ValidationResult {
	if r.URL == "" {
		return ValidationResult{Err: errors.New("empty url")}
	}
	parsedURL, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return ValidationResult{Err: fmt.Errorf("invalid url: %w", err)}
	}
	return ValidationResult{Domain: parsedURL.Host, Err: nil}
}

type ShortenResponse struct {
	URL string `json:"url"`
}
