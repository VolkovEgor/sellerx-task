package service

import "github.com/VolkovEgor/sellerx-task/internal/repository"

type User interface {
}

type Chat interface {
}

type Message interface {
}

type Service struct {
	User
	Chat
	Message
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
