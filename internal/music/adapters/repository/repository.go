package repository

import (
	"gorm.io/gorm"
)

type MusicRepository struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) *MusicRepository {
	// db.AutoMigrate(&model.Artist{}, &model.Genre{}, &model.Song{})
	return &MusicRepository{
		db: db,
	}
}

func (r *MusicRepository) IMockRepository() {}
