package playback

import (
	"music-app-backend/internal/playback/adapters/repository"
	"music-app-backend/internal/playback/application"
	"music-app-backend/internal/playback/ports"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaybackModule struct {
	Repository ports.IPlaybackRepository
	Service    ports.IPlaybackService
}

func NewPlaybackModule(db *gorm.DB) *PlaybackModule {
	playbackRepo := repository.NewPlaybackRepository(db)
	playbackService := application.NewPlaybackService(playbackRepo)

	return &PlaybackModule{
		Repository: playbackRepo,
		Service:    playbackService,
	}
}

func (s *PlaybackModule) RegisterRoutes(router *gin.Engine) {}
