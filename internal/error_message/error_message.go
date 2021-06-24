package error_message

import "errors"

var (
	ErrWrongUsername     = errors.New("username must contain from 1 to 50 characters")
	ErrWrongChatname     = errors.New("chat name must contain from 1 to 50 characters")
	ErrNoChatUsers       = errors.New("chat must contain least one user")
	ErrChatUserNotExists = errors.New("chat user not exists")
)
