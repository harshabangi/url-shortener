package util

import (
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func TestMd5ToBase62(t *testing.T) {
	assert := asserts.New(t)

	tcc := []struct {
		name      string
		md5Hash   string
		want      string
		shouldErr bool
	}{
		{"Valid md5 hash", "5eb63bbbe01eeed093cb22bb8f5acdc3", "c2SIRhqtOH9heLczsBsHm3", false},
		{"Invalid MD5 Hash", "invalid-md5-hash", "", true},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Md5ToBase62(tc.md5Hash)
			if tc.shouldErr {
				assert.NotNil(err)
			} else {
				assert.Equal(tc.want, got)
			}
		})
	}
}
