package user

import (
	"music-app-backend/internal/user/adapters/repository"
	"music-app-backend/internal/user/application"
	"music-app-backend/internal/user/ports"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserModule struct {
	Repository ports.IUserRepository
	Service    ports.IUserService
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepo := repository.NewUserRepository(db)
	userService := application.NewUserService(userRepo)

	return &UserModule{
		Repository: userRepo,
		Service:    userService,
	}
}

func (u *UserModule) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/users")
}
