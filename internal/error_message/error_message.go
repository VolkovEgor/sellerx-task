package error_message

import "errors"

var (
	ErrWrongUsername      = errors.New("username must contain from 1 to 50 characters")
	ErrWrongChatname      = errors.New("chat name must contain from 1 to 50 characters")
	ErrWrongTextMes       = errors.New("message must not be empty")
	ErrNoChatUsers        = errors.New("chat must contain least one user")
	ErrUserNotExists      = errors.New("user not exists")
	ErrMesAuthorNotExists = errors.New("message author not exists")
	ErrChatNotExists      = errors.New("chat not exists")
)
