package model

import (
	"music-app-backend/pkg/model"
)

type ArtistFollower struct {
	model.BaseModel
	ArtistID            uint64 `json:"artist_id" gorm:"not null;uniqueIndex:idx_artist_follower"`
	FollowerUserID      uint64 `json:"follower_user_id" gorm:"not null;uniqueIndex:idx_artist_follower"`
	NotificationEnabled bool   `json:"notification_enabled" gorm:"default:true"`
	FollowedAt          int64  `json:"followed_at"`
}
