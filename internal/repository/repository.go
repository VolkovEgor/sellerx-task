package repository

import (
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user *model.User) (string, error)
	ExistenceCheck(userId string) error
}

type Chat interface {
	Create(user *model.Chat) (string, error)
	GetAllForUser(userId string) ([]*model.Chat, error)
	ExistenceCheck(chatId string) error
}

type Message interface {
	Create(message *model.Message) (string, error)
	GetAllForChat(chatId string) ([]*model.Message, error)
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
