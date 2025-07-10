package repository

import (
	model "music-app-backend/internal/music/domain"
	"music-app-backend/pkg/database"

	"gorm.io/gorm"
)

type MusicRepository struct {
	database.Repository
}

func NewMusicRepository(db *gorm.DB) *MusicRepository {
	db.AutoMigrate(&model.Artist{}, &model.Genre{}, &model.Song{})
	return &MusicRepository{
		Repository: database.NewRepository(db),
	}
}
