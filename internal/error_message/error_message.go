package error_message

import "errors"

var (
	ErrWrongUsername        = errors.New("username must contain from 1 to 50 characters")
	ErrWrongChatname        = errors.New("chat name must contain from 1 to 50 characters")
	ErrWrongMesText         = errors.New("message must not be empty")
	ErrNoChatUsers          = errors.New("chat must contain least one user")
	ErrRecurringUsers       = errors.New("chat users must not repeat themselves")
	ErrWrongMesCreationTime = errors.New("message time must be longer than chat time")
	ErrUserNotExists        = errors.New("user not exists")
	ErrChatNotExists        = errors.New("chat not exists")
	ErrEmptyUserId          = errors.New("empty user id")
	ErrEmptyChatId          = errors.New("empty chat id")
)
