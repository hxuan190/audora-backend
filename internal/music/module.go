package music

import (
	"music-app-backend/internal/music/adapters/repository"
	"music-app-backend/internal/music/application"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MusicModule struct {
	Repository *repository.MusicRepository
	Service    *application.MusicService
}

func NewMusicModule(db *gorm.DB) *MusicModule {
	musicRepo := repository.NewMusicRepository(db)
	musicService := application.NewMusicService(musicRepo)

	return &MusicModule{
		Repository: musicRepo,
		Service:    musicService,
	}
}

func (s *MusicModule) RegisterRoutes(router *gin.RouterGroup) {}
