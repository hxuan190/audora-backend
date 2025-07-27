package application

import (
	"music-app-backend/internal/user/adapters/repository"
	model "music-app-backend/internal/user/domain"
	baseModel "music-app-backend/pkg/model"

	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo  *repository.UserRepository
	generator *goflakeid.Generator
}

func NewUserService(userRepo *repository.UserRepository, generator *goflakeid.Generator) *UserService {
	return &UserService{
		userRepo:  userRepo,
		generator: generator,
	}
}

func (s *UserService) CreateUserAfterRegistration(user *model.AfterRegistrationRequest) (*uint64, error) {
	identityID, err := uuid.Parse(user.Identity.ID)
	if err != nil {
		return nil, err
	}

	baseModelInstance, err := baseModel.NewBaseModel(s.generator)
	if err != nil {
		return nil, err
	}

	userModel, err := s.userRepo.CreateUserAfterRegistration(&model.User{
		BaseModel:        *baseModelInstance,
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
