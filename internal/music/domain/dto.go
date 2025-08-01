package model

type CreateArtistDTO struct {
	UserID          uint64  `json:"user_id" gorm:"not null;unique"`
	ArtistName      string  `json:"artist_name" gorm:"not null;size:150"`
	Bio             *string `json:"bio"`
	ProfileImageURL *string `json:"profile_image_url"`
	BannerImageURL  *string `json:"banner_image_url"`
	WebsiteURL      *string `json:"website_url"`
	SpotifyURL      *string `json:"spotify_url"`
	InstagramURL    *string `json:"instagram_url"`
	TwitterURL      *string `json:"twitter_url"`
	YoutubeURL      *string `json:"youtube_url"`
}
