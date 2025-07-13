package repository

import (
	"gorm.io/gorm"
)

type PlaybackRepository struct {
	db *gorm.DB
}

func NewPlaybackRepository(db *gorm.DB) *PlaybackRepository {
	// db.AutoMigrate(&model.ListeningSession{}, &model.Mood{}, &model.PlayListSongs{}, &model.PlayList{}, &model.SongPlay{})
	return &PlaybackRepository{
		db: db,
	}
}

func (r *PlaybackRepository) IMockRepository() {}
