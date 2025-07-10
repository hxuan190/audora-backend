package application

import (
	model "music-app-backend/internal/user/domain"
	"music-app-backend/internal/user/ports"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo ports.IUserRepository
}

func NewUserService(userRepo ports.IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUserAfterRegistration(user *model.AfterRegistrationRequest) (*uuid.UUID, error) {
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
