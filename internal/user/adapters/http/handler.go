package http

import (
	model "music-app-backend/internal/user/domain"
	"music-app-backend/internal/user/ports"
	app_error "music-app-backend/pkg/error"
	json_response "music-app-backend/pkg/json"

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

func (h *UserHandler) HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if appErr, ok := app_error.GetAppError(err); ok {
		json_response.ResponseJSON(c, appErr.StatusCode, appErr.Message, appErr.Data)
		return true
	}

	json_response.ResponseInternalError(c, err)
	return true
}

func (h *UserHandler) AfterRegistration(c *gin.Context) {
	var request model.AfterRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		json_response.ResponseBadRequest(c, err.Error())
		return
	}

	userID, err := h.userService.CreateUserAfterRegistration(&request)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	json_response.ResponseOK(c, userID)
}
