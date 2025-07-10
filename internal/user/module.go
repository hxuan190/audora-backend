package user

import (
	"music-app-backend/internal/user/adapters/http"
	"music-app-backend/internal/user/adapters/repository"
	"music-app-backend/internal/user/application"
	"music-app-backend/internal/user/ports"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserModule struct {
	Repository ports.IUserRepository
	Service    ports.IUserService
	Handler    *http.UserHandler
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepo := repository.NewUserRepository(db)
	userService := application.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	return &UserModule{
		Repository: userRepo,
		Service:    userService,
		Handler:    userHandler,
	}
}

func (u *UserModule) RegisterRoutes(router *gin.RouterGroup) {
	internal := router.Group("internal")
	internal.POST("/hooks/after-registration", u.Handler.AfterRegistration)
}
