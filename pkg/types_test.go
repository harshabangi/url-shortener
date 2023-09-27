package pkg

import (
	"errors"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func TestShortenRequestValidation(t *testing.T) {
	tcc := []struct {
		name string
		req  ShortenRequest
		want ValidationResult
	}{
		{"Valid URL", ShortenRequest{URL: "https://example.com"}, ValidationResult{Domain: "example.com"}},
		{"Empty URL", ShortenRequest{URL: ""}, ValidationResult{Err: errors.New("empty url")}},
		{"Invalid URL", ShortenRequest{URL: "invalid-url"}, ValidationResult{Err: errors.New("invalid url")}},
	}

	assert := asserts.New(t)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.req.Validate()
			if tc.want.Err != nil {
				assert.NotNil(got.Err)
			} else {
				assert.Equal(tc.want, got)
			}
		})
	}
}
