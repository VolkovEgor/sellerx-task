package service

import (
	"time"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user *model.User) (string, error) {
	if user.Username == "" || len(user.Username) > 50 {
		return "", errMes.ErrWrongUsername
	}

	user.CreatedAt = time.Now().Unix()
	return s.repo.Create(user)
}
