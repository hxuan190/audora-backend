package music

import (
	"music-app-backend/internal/music/adapters/http"
	"music-app-backend/internal/music/adapters/repository"
	"music-app-backend/internal/music/application"
	ctx2 "music-app-backend/pkg/context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MusicModule struct {
	Repository repository.IMusicRepository
	Service    *application.MusicService
	Handler    *http.MusicHandler
}

func NewMusicModule(db *gorm.DB, serviceContext *ctx2.ServiceContext) *MusicModule {
	musicRepo := repository.NewMusicRepository(db)
	musicService := application.NewMusicService(musicRepo, serviceContext.GetIDGenerator())
	uploadHandler := http.NewMusicHandler(musicService, serviceContext.GetStorageService(), serviceContext.GetRedisClient(), serviceContext.GetIDGenerator())

	return &MusicModule{
		Repository: musicRepo,
		Service:    musicService,
		Handler:    uploadHandler,
	}
}

func (s *MusicModule) RegisterRoutes(router *gin.RouterGroup) {
	uploadRouter := router.Group("/upload")
	{
		uploadRouter.POST("/initiate", s.Handler.InitiateUpload)
		uploadRouter.POST("/complete", s.Handler.CompleteUpload)
		uploadRouter.GET("/status/{upload_id}", s.Handler.GetUploadStatus)
	}
	router.GET("/processing/status/{song_id}", s.Handler.GetProcessingStatus)
	router.POST("/processing/callback/{song_id}", s.Handler.ProcessingCallback)
	streamRouter := router.Group("/stream")
	{
		streamRouter.GET("/{song_id}", s.Handler.GetStreamingURL)
	}
}
