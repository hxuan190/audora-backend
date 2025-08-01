package music

import (
	ctx2 "music-app-backend/pkg/context"
	"music-app-backend/internal/music/adapters/repository"
	"music-app-backend/internal/music/application"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MusicModule struct {
	Repository repository.IMusicRepository
	Service    *application.MusicService
}

func NewMusicModule(db *gorm.DB, serviceContext *ctx2.ServiceContext) *MusicModule {
	musicRepo := repository.NewMusicRepository(db)
	musicService := application.NewMusicService(musicRepo, serviceContext.GetIDGenerator())

	return &MusicModule{
		Repository: musicRepo,
		Service:    musicService,
	}
}

func (s *MusicModule) RegisterRoutes(router *gin.RouterGroup) {}
