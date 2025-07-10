package repository

import (
	model "music-app-backend/internal/analytics/domain"
	"music-app-backend/pkg/database"

	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	database.Repository
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	db.AutoMigrate(&model.DailyArtistStats{}, &model.DailySongStats{})
	return &AnalyticsRepository{
		Repository: database.NewRepository(db),
	}
}
