package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type UserPreference struct {
	model.BaseModel
	UserID                     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;unique"`
	PreferredGenres            []int     `json:"preferred_genres" gorm:"type:integer[]"`
	PreferredMoods             []int     `json:"preferred_moods" gorm:"type:integer[]"`
	AutoPlay                   bool      `json:"auto_play" gorm:"default:true"`
	ShuffleByDefault           bool      `json:"shuffle_by_default" gorm:"default:false"`
	NotificationNewReleases    bool      `json:"notification_new_releases" gorm:"default:true"`
	NotificationArtistMessages bool      `json:"notification_artist_messages" gorm:"default:true"`
	NotificationTipsReceived   bool      `json:"notification_tips_received" gorm:"default:true"`
	ExplicitContentAllowed     bool      `json:"explicit_content_allowed" gorm:"default:false"`
}
