package http

import (
	"music-app-backend/internal/music/application"
	model "music-app-backend/internal/music/domain"
	app_error "music-app-backend/pkg/error"
	json_response "music-app-backend/pkg/json"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	musicService *application.MusicService
}

func NewMusicHandler(musicService *application.MusicService) *MusicHandler {
	return &MusicHandler{
		musicService: musicService,
	}
}

func (h *MusicHandler) HandleError(c *gin.Context, err error) bool {
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

func (h *MusicHandler) InitiateUpload(c *gin.Context) {
	uploadRequest := &model.InitiateUploadRequest{}
	if err := c.ShouldBindJSON(uploadRequest); err != nil {
		json_response.ResponseBadRequest(c, "Invalid request data")
		return
	}

	response, err := h.musicService.InitiateUpload(c.Request.Context(), uploadRequest)
	if h.HandleError(c, err) {
		return
	}

	json_response.ResponseJSON(c, 200, "Upload initiated successfully", response)
}
