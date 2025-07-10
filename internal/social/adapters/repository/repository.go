package repository

import (
	model "music-app-backend/internal/social/domain"

	"gorm.io/gorm"
)

type SocialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) *SocialRepository {
	db.AutoMigrate(&model.ArtistFollower{}, &model.ArtistFollower{}, &model.ArtistMessage{})
	return &SocialRepository{
		db: db,
	}
}

func (r *SocialRepository) IMockRepository() {}
