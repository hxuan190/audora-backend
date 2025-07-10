package music

import (
	"music-app-backend/internal/music/adapters/repository"
	"music-app-backend/internal/music/application"
	"music-app-backend/internal/music/ports"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MusicModule struct {
	Repository ports.IMusicRepository
	Service    ports.IMusicService
}

func NewMusicModule(db *gorm.DB) *MusicModule {
	musicRepo := repository.NewMusicRepository(db)
	musicService := application.NewMusicService(musicRepo)

	return &MusicModule{
		Repository: musicRepo,
		Service:    musicService,
	}
}

func (s *MusicModule) RegisterRoutes(router *gin.Engine) {}
