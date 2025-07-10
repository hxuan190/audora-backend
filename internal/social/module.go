package social

import (
	"music-app-backend/internal/social/adapters/repository"
	"music-app-backend/internal/social/application"
	"music-app-backend/internal/social/ports"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialModule struct {
	Repository ports.ISocialRepository
	Service    ports.ISocialService
}

func NewSocialModule(db *gorm.DB) *SocialModule {
	socialRepo := repository.NewSocialRepository(db)
	socialService := application.NewSocialService(socialRepo)

	return &SocialModule{
		Repository: socialRepo,
		Service:    socialService,
	}
}

func (s *SocialModule) RegisterRoutes(router *gin.RouterGroup) {}
