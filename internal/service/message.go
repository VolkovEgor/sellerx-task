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
		return "", errMes.ErrWrongMesText
	}

	if message.ChatId == "" {
		return "", errMes.ErrEmptyChatId
	}

	if message.AuthorId == "" {
		return "", errMes.ErrEmptyUserId
	}

	chat, err := s.chatRepo.GetById(message.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errMes.ErrChatNotExists
		}
		return "", err
	}

	if err := s.userRepo.ExistenceCheck(message.AuthorId); err != nil {
		if err == sql.ErrNoRows {
			return "", errMes.ErrUserNotExists
		}
		return "", err
	}

	message.CreatedAt = time.Now().Unix()
	if message.CreatedAt > chat.CreatedAt {
		return "", errMes.ErrWrongMesCreationTime
	}

	return s.repo.Create(message)
}

func (s *MessageService) GetAllForChat(chatId string) ([]*model.Message, error) {
	if chatId == "" {
		return nil, errMes.ErrEmptyChatId
	}

	if err := s.chatRepo.ExistenceCheck(chatId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errMes.ErrChatNotExists
		}
		return nil, err
	}

	return s.repo.GetAllForChat(chatId)
}
