package analytics

import (
	"music-app-backend/internal/analytics/adapters/repository"
	"music-app-backend/internal/analytics/application"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsModule struct {
	Repository *repository.AnalyticsRepository
	Service    *application.AnalyticsService
}

func NewAnalyticsModule(db *gorm.DB) *AnalyticsModule {
	analyticsRepo := repository.NewAnalyticsRepository(db)
	analyticsService := application.NewAnalyticsService(analyticsRepo)

	return &AnalyticsModule{
		Repository: analyticsRepo,
		Service:    analyticsService,
	}
}

func (s *AnalyticsModule) RegisterRoutes(router *gin.RouterGroup) {}
