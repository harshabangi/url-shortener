package memory

import (
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTopNDomainsByFrequency(t *testing.T) {
	tcc := []struct {
		name  string
		input map[string]int64
		limit int
		want  []shared.DomainFrequency
	}{
		{"test empty input map", map[string]int64{}, 3, nil},
		{
			"test non empty input map with limit 3",
			map[string]int64{
				"example1.com": 10,
				"example2.com": 5,
				"example3.com": 15,
				"example4.com": 8,
				"example5.com": 12,
			},
			3,
			[]shared.DomainFrequency{
				{Domain: "example3.com", Frequency: 15},
				{Domain: "example5.com", Frequency: 12},
				{Domain: "example1.com", Frequency: 10},
			},
		},
		{
			"test non empty input map with limit 3 and entries less than 3",
			map[string]int64{
				"example3.com": 15,
				"example5.com": 12,
			},
			3,
			[]shared.DomainFrequency{
				{Domain: "example3.com", Frequency: 15},
				{Domain: "example5.com", Frequency: 12},
			},
		},
	}

	assert := asserts.New(t)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			got := getTopNDomainsByFrequency(tc.input, tc.limit)
			assert.Equal(tc.want, got)
		})
	}
}
