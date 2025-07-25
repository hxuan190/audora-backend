package model

import (
	"music-app-backend/pkg/model"
)

type ListeningSession struct {
	model.BaseModel
	UserID        *uint64 `json:"user_id"`
	SessionID     string  `json:"session_id" gorm:"not null;size:100"`
	SongID        uint64  `json:"song_id" gorm:"not null"`
	ArtistID      uint64  `json:"artist_id" gorm:"not null;index"`
	CountryCode   string  `json:"country_code" gorm:"size:2"`
	City          string  `json:"city" gorm:"size:100"`
	StartedAt     int64   `json:"started_at"`
	LastHeartbeat int64   `json:"last_heartbeat"`
	IsActive      bool    `json:"is_active" gorm:"default:true;index:idx_active_heartbeat"`
}
