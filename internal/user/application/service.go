package application

import (
	"music-app-backend/internal/user/ports"
)

type UserService struct {
	userRepo ports.IUserRepository
}

func NewUserService(userRepo ports.IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) IMockService() {}
