package service

import (
	"database/sql"
	"time"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository"
)

type ChatService struct {
	repo     repository.Chat
	userRepo repository.User
}

func NewChatService(repo repository.Chat, userRepo repository.User) *ChatService {
	return &ChatService{repo: repo, userRepo: userRepo}
}

func (s *ChatService) Create(chat *model.Chat) (string, error) {
	if chat.Name == "" || len(chat.Name) > 50 {
		return "", errMes.ErrWrongChatname
	}

	if len(chat.Users) == 0 {
		return "", errMes.ErrNoChatUsers
	}

	set := make(map[string]bool)
	for _, userId := range chat.Users {
		if userId == "" {
			return "", errMes.ErrEmptyUserId
		}

		if set[userId] {
			return "", errMes.ErrRecurringUsers
		}
		set[userId] = true

		if err := s.userRepo.ExistenceCheck(userId); err != nil {
			if err == sql.ErrNoRows {
				return "", errMes.ErrUserNotExists
			}
			return "", err
		}
	}

	chat.CreatedAt = time.Now().Unix()
	return s.repo.Create(chat)
}

func (s *ChatService) GetAllForUser(userId string) ([]*model.Chat, error) {
	if userId == "" {
		return nil, errMes.ErrEmptyUserId
	}

	if err := s.userRepo.ExistenceCheck(userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errMes.ErrUserNotExists
		}
		return nil, err
	}

	return s.repo.GetAllForUser(userId)
}
