package repository

import (
	model "music-app-backend/internal/playback/domain"
	"music-app-backend/pkg/database"

	"gorm.io/gorm"
)

type PlaybackRepository struct {
	database.Repository
}

func NewPlaybackRepository(db *gorm.DB) *PlaybackRepository {
	db.AutoMigrate(&model.ListeningSession{}, &model.Mood{}, &model.PlayListSongs{}, &model.PlayList{}, &model.SongPlay{})
	return &PlaybackRepository{
		Repository: database.NewRepository(db),
	}
}
