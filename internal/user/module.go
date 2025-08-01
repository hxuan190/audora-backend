package user

import (
	musicModuleSvc "music-app-backend/internal/music/application"
	"music-app-backend/internal/user/adapters/http"
	"music-app-backend/internal/user/adapters/repository"
	"music-app-backend/internal/user/application"
	ctx2 "music-app-backend/pkg/context"

	"github.com/gin-gonic/gin"
)

type UserModule struct {
	Repository *repository.UserRepository
	Service    *application.UserService
	Handler    *http.UserHandler

	artistService musicModuleSvc.IMusicService
}

func NewUserModule(serviceContext *ctx2.ServiceContext, artistService musicModuleSvc.IMusicService) *UserModule {
	userRepo := repository.NewUserRepository(serviceContext.GetDB())
	userService := application.NewUserService(userRepo, serviceContext.GetIDGenerator(), artistService)
	userHandler := http.NewUserHandler(userService)

	return &UserModule{
		Repository: userRepo,
		Service:    userService,
		Handler:    userHandler,
		artistService: artistService,
	}
}

func (u *UserModule) RegisterRoutes(router *gin.RouterGroup) {
	internal := router.Group("internal")
	internal.POST("/hooks/after-registration", u.Handler.AfterRegistration)
}
