package ports

import (
	model "music-app-backend/internal/user/domain"

	"github.com/google/uuid"
)

type IUserRepository interface {
	CreateUserAfterRegistration(user *model.User) (*model.User, error)
}

type IUserService interface {
	CreateUserAfterRegistration(user *model.AfterRegistrationRequest) (*uuid.UUID, error)
}
