package repository

import (
	"gorm.io/gorm"
)

type SocialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) *SocialRepository {
	// db.AutoMigrate(&model.ArtistFollower{}, &model.ArtistFollower{}, &model.ArtistMessage{})
	return &SocialRepository{
		db: db,
	}
}

func (r *SocialRepository) IMockRepository() {}
