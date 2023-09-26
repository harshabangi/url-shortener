package pkg

import (
	"fmt"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func (s *ShortenRequest) ValidateURL() error {
	_, err := url.ParseRequestURI(s.URL)
	if err == nil {
		return nil
	}
	return fmt.Errorf("error validating URL: %s: %w", s.URL, err)
}

type ShortenResponse struct {
	URL string `json:"url"`
}
