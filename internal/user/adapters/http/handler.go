package http

import (
	"music-app-backend/internal/user/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService ports.IUserService
}

func NewUserHandler(userService ports.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) AfterRegistration(c *gin.Context) {

}
