package service

import (
	"database/sql"
	"time"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository"
)

type MessageService struct {
	repo     repository.Message
	userRepo repository.User
	chatRepo repository.Chat
}

func NewMessageService(repo repository.Message, userRepo repository.User, chatRepo repository.Chat) *MessageService {
	return &MessageService{
		repo:     repo,
		userRepo: userRepo,
		chatRepo: chatRepo,
	}
}

func (s *MessageService) Create(message *model.Message) (string, error) {
	if message.Text == "" {
		return "", errMes.ErrWrongTextMes
	}

	if message.ChatId == "" {
		return "", errMes.ErrChatNotExists
	}

	if message.AuthorId == "" {
		return "", errMes.ErrMesAuthorNotExists
	}

	if err := s.chatRepo.ExistenceCheck(message.ChatId); err != nil {
		if err == sql.ErrNoRows {
			return "", errMes.ErrChatNotExists
		}
		return "", err
	}

	if err := s.userRepo.ExistenceCheck(message.AuthorId); err != nil {
		if err == sql.ErrNoRows {
			return "", errMes.ErrMesAuthorNotExists
		}
		return "", err
	}

	message.CreatedAt = time.Now().Unix()
	return s.repo.Create(message)
}

func (s *MessageService) GetAllForChat(chatId string) ([]*model.Message, error) {
	if chatId == "" {
		return nil, errMes.ErrChatNotExists
	}

	if err := s.chatRepo.ExistenceCheck(chatId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errMes.ErrChatNotExists
		}
		return nil, err
	}

	return s.repo.GetAllForChat(chatId)
}
