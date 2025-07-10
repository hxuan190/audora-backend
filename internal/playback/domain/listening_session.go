package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type ListeningSession struct {
	model.BaseModel
	UserID        *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	SessionID     string     `json:"session_id" gorm:"not null;size:100"`
	SongID        uuid.UUID  `json:"song_id" gorm:"type:uuid;not null"`
	ArtistID      uuid.UUID  `json:"artist_id" gorm:"type:uuid;not null;index"`
	CountryCode   string     `json:"country_code" gorm:"size:2"`
	City          string     `json:"city" gorm:"size:100"`
	StartedAt     int64      `json:"started_at"`
	LastHeartbeat int64      `json:"last_heartbeat"`
	IsActive      bool       `json:"is_active" gorm:"default:true;index:idx_active_heartbeat"`
}
