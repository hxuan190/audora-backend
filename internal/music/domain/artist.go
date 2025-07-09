package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type Artist struct {
	model.BaseModel
	UserID                uuid.UUID `json:"user_id" gorm:"unique"`
	ArtistName            string    `json:"artist_name"`
	ArtistBio             string    `json:"artist_bio"`
	ProfileImageURL       string    `json:"profile_image_url"`
	BannerImageURL        string    `json:"banner_image_url"`
	WebsiteURL            string    `json:"website_url"`
	SpotifyURL            string    `json:"spotify_url"`
	InstagramURL          string    `json:"instagram_url"`
	TwitterURL            string    `json:"twitter_url"`
	YoutubeURL            string    `json:"youtube_url"`
	IsVerified            bool      `json:"is_verified" gorm:"default:false"`
	VerificationRequestAt int64     `json:"verification_request_at" gorm:"autoUpdateTime"`
	TotalPlays            int64     `json:"total_plays" gorm:"default:0"`
	TotalEarnings         int64     `json:"total_earnings" gorm:"default:0"`
	FlowerCount           int64     `json:"flower_count" gorm:"default:0"`
	IsActive              bool      `json:"is_active" gorm:"default:true"`
}
