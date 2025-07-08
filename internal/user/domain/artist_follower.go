package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type ArtistFollower struct {
	model.BaseModel
	ArtistID            uuid.UUID `json:"artist_id" gorm:"not null;uniqueIndex:idx_artist_follower"`
	FollowerUserID      uuid.UUID `json:"follower_user_id" gorm:"not null;uniqueIndex:idx_artist_follower"`
	NotificationEnabled bool      `json:"notification_enabled" gorm:"default:true"`
	FollowedAt          int64     `json:"followed_at" gorm:"default:CURRENT_TIMESTAMP"`
}
