package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type DailyArtistStats struct {
	model.BaseModel
	ArtistID            uuid.UUID `json:"artist_id" gorm:"not null;uniqueIndex:idx_artist_date"`
	Date                string    `json:"date" gorm:"not null;uniqueIndex:idx_artist_date"`
	TotalPlays          int       `json:"total_plays" gorm:"default:0"`
	UniqueListeners     int       `json:"unique_listeners" gorm:"default:0"`
	TotalDurationPlayed int       `json:"total_duration_played" gorm:"default:0"`
	NewFollowers        int       `json:"new_followers" gorm:"default:0"`
	TipsReceivedCents   int       `json:"tips_received_cents" gorm:"default:0"`
	TipCount            int       `json:"tip_count" gorm:"default:0"`
}
