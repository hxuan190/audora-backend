package model

import (
	"music-app-backend/pkg/model"
)

type Artist struct {
	model.BaseModel
	UserID                  uint64  `json:"user_id" gorm:"not null;unique"`
	ArtistName              string  `json:"artist_name" gorm:"not null;size:150"`
	Bio                     string  `json:"bio"`
	ProfileImageURL         string  `json:"profile_image_url"`
	BannerImageURL          string  `json:"banner_image_url"`
	WebsiteURL              string  `json:"website_url"`
	SpotifyURL              string  `json:"spotify_url"`
	InstagramURL            string  `json:"instagram_url"`
	TwitterURL              string  `json:"twitter_url"`
	YoutubeURL              string  `json:"youtube_url"`
	IsVerified              bool    `json:"is_verified" gorm:"default:false"`
	VerificationRequestedAt *int64  `json:"verification_requested_at"`
	TotalPlays              int64   `json:"total_plays" gorm:"default:0"`
	TotalEarnings           float64 `json:"total_earnings" gorm:"type:decimal(10,2);default:0.00"`
	FollowerCount           int     `json:"follower_count" gorm:"default:0"`
}
