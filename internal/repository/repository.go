package repository

import (
	"github.com/jmoiron/sqlx"
)

type User interface {
}

type Chat interface {
}

type Message interface {
}

type Repository struct {
	User
	Chat
	Message
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
