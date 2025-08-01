package repository

import (
	model "music-app-backend/internal/music/domain"

	"gorm.io/gorm"
)

type IMusicRepository interface {
	InsertArtist(artist *model.Artist) error
}

type MusicRepository struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) *MusicRepository {
	return &MusicRepository{
		db: db,
	}
}

func (db *MusicRepository) InsertArtist(artist *model.Artist) error {
	return db.db.Create(artist).Error
}
