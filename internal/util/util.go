package util

import (
	"encoding/hex"
	"fmt"
	"github.com/eknkc/basex"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Md5ToBase62(md5Hash string) (string, error) {
	md5Bytes, err := hex.DecodeString(md5Hash)
	if err != nil {
		return "", fmt.Errorf("error decoding MD5 hash: %w", err)
	}

	encoder, err := basex.NewEncoding(base62Chars)
	if err != nil {
		return "", fmt.Errorf("error initializing Base62 encoder: %w", err)
	}

	return encoder.Encode(md5Bytes), nil
}
