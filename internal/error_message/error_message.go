package error_message

import "errors"

var (
	ErrWrongUsername = errors.New("username must contain from 1 to 50 characters")
)
