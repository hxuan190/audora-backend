package application

import (
	"music-app-backend/internal/user/adapters/repository"
	model "music-app-backend/internal/user/domain"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUserAfterRegistration(user *model.AfterRegistrationRequest) (*uint64, error) {
	identityID, err := uuid.Parse(user.Identity.ID)
	if err != nil {
		return nil, err
	}

	userModel, err := s.userRepo.CreateUserAfterRegistration(&model.User{
		KratosIdentityID: identityID,
		Email:            user.Identity.Traits.Email,
		UserType:         user.Identity.Traits.UserType,
		DisplayName:      user.Identity.Traits.DisplayName,
		AvatarURL:        "",
		IsActive:         true,
		LastLoginAt:      nil,
	})

	if err != nil {
		return nil, err
	}

	return &userModel.ID, nil
}
