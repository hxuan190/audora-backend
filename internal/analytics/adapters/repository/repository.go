package repository

import (
	model "music-app-backend/internal/analytics/domain"

	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	db.AutoMigrate(&model.DailyArtistStats{}, &model.DailySongStats{})
	return &AnalyticsRepository{
		db: db,
	}
}

func (r *AnalyticsRepository) IMockRepository() {}
