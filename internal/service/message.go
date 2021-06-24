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

func (s *MessageService) Create(message *model.Message) (int, error) {
	if message.Text == "" {
		return 0, errMes.ErrWrongTextMes
	}

	if message.ChatId <= 0 {
		return 0, errMes.ErrChatNotExists
	}

	if message.AuthorId <= 0 {
		return 0, errMes.ErrMesAuthorNotExists
	}

	if err := s.chatRepo.ExistenceCheck(message.ChatId); err != nil {
		if err == sql.ErrNoRows {
			return 0, errMes.ErrChatNotExists
		}
		return 0, err
	}

	if err := s.userRepo.ExistenceCheck(message.AuthorId); err != nil {
		if err == sql.ErrNoRows {
			return 0, errMes.ErrMesAuthorNotExists
		}
		return 0, err
	}

	message.CreatedAt = time.Now().Unix()
	return s.repo.Create(message)
}
