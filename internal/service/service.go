package service

import (
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository"
)

type User interface {
	Create(user *model.User) (int, error)
}

type Chat interface {
	Create(user *model.Chat) (int, error)
}

type Message interface {
}

type Service struct {
	User
	Chat
	Message
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
		Chat: NewChatService(repos.Chat, repos.User),
	}
}
