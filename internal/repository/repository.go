package repository

import (
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user *model.User) (int, error)
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
	return &Repository{
		User: postgres.NewUserPg(db),
	}
}
