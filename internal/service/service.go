package service

import (
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository"
)

type User interface {
	Create(user *model.User) (string, error)
}

type Chat interface {
	Create(user *model.Chat) (string, error)
	GetAllForUser(userId string) ([]*model.Chat, error)
}

type Message interface {
	Create(message *model.Message) (string, error)
	GetAllForChat(chatId string) ([]*model.Message, error)
}

type Service struct {
	User
	Chat
	Message
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:    NewUserService(repos.User),
		Chat:    NewChatService(repos.Chat, repos.User),
		Message: NewMessageService(repos.Message, repos.User, repos.Chat),
	}
}
