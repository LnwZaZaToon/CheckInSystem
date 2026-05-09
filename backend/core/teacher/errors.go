package teacher

import "errors"

var (
	ErrNotFound       = errors.New("teacher not found")
	ErrUsernameExists = errors.New("username already exists")
)
