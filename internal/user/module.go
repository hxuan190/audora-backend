package user

import (
	ctx2 "music-app-backend/pkg/context"
	"music-app-backend/internal/user/adapters/http"
	"music-app-backend/internal/user/adapters/repository"
	"music-app-backend/internal/user/application"

	"github.com/gin-gonic/gin"
)

type UserModule struct {
	Repository *repository.UserRepository
	Service    *application.UserService
	Handler    *http.UserHandler
}

func NewUserModule(serviceContext *ctx2.ServiceContext) *UserModule {
	userRepo := repository.NewUserRepository(serviceContext.GetDB())
	userService := application.NewUserService(userRepo, serviceContext.GetIDGenerator())
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
