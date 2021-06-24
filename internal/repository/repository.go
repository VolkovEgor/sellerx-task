package repository

import (
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user *model.User) (int, error)
	ExistenceCheck(userId int) error
}

type Chat interface {
	Create(user *model.Chat) (int, error)
	ExistenceCheck(chatId int) error
}

type Message interface {
	Create(message *model.Message) (int, error)
}

type Repository struct {
	User
	Chat
	Message
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:    postgres.NewUserPg(db),
		Chat:    postgres.NewChatPg(db),
		Message: postgres.NewMessagePg(db),
	}
}
