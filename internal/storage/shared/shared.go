package shared

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrCollision = errors.New("item already exists")
)
