package repository

import (
	model "music-app-backend/internal/social/domain"
	"music-app-backend/pkg/database"

	"gorm.io/gorm"
)

type SocialRepository struct {
	database.Repository
}

func NewSocialRepository(db *gorm.DB) *SocialRepository {
	db.AutoMigrate(&model.ArtistFollower{}, &model.ArtistFollower{}, &model.ArtistMessage{})
	return &SocialRepository{
		Repository: database.NewRepository(db),
	}
}
