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

func (s *ChatService) Create(chat *model.Chat) (int, error) {
	if chat.Name == "" || len(chat.Name) > 50 {
		return 0, errMes.ErrWrongChatname
	}

	if len(chat.Users) == 0 {
		return 0, errMes.ErrNoChatUsers
	}

	for _, userId := range chat.Users {
		if err := s.userRepo.ExistenceCheck(userId); err != nil {
			if err == sql.ErrNoRows {
				return 0, errMes.ErrChatUserNotExists
			}
			return 0, err
		}
	}

	chat.CreatedAt = time.Now().Unix()
	return s.repo.Create(chat)
}
