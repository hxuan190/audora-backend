package user

import (
	"music-app-backend/internal/user/adapters/repository"
	"music-app-backend/internal/user/application"
	"music-app-backend/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserModule struct {
	DB         *database.Database
	Repository *repository.UserRepository
	Service    *application.UserService
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepo := repository.NewUserRepository(db)
	userService := application.NewUserService(userRepo)

	return &UserModule{
		Repository: userRepo,
		Service:    userService,
	}
}

func (u *UserModule) RegisterRoutes(router *gin.Engine) {
	router.POST("/users")
}
