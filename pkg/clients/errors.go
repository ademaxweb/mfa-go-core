package clients

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrInvalidData        = errors.New("invalid user data")
	ErrServiceUnavailable = errors.New("service unavailable")
)
