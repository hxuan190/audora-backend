package analytics

import (
	"music-app-backend/internal/analytics/adapters/repository"
	"music-app-backend/internal/analytics/application"
	"music-app-backend/internal/analytics/ports"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsModule struct {
	Repository ports.IAnalyticsRepository
	Service    ports.IAnalyticsService
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
