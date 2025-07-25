package model

import (
	"music-app-backend/pkg/model"
)

type SkipReason string

const (
	SkipReasonUserSkip      SkipReason = "user_skip"
	SkipReasonNextSong      SkipReason = "next_song"
	SkipReasonEndOfPlaylist SkipReason = "end_of_playlist"
	SkipReasonError         SkipReason = "error"
)

type SongPlay struct {
	model.BaseModel
	SongID                uint64      `json:"song_id" gorm:"not null;index"`
	UserID                *uint64     `json:"user_id" gorm:"index"`
	SessionID             string      `json:"session_id" gorm:"size:100"`
	IPAddress             string      `json:"ip_address"`
	UserAgent             string      `json:"user_agent"`
	CountryCode           string      `json:"country_code" gorm:"size:2"`
	City                  string      `json:"city" gorm:"size:100"`
	DurationPlayedSeconds int         `json:"duration_played_seconds" gorm:"default:0"`
	Completed             bool        `json:"completed" gorm:"default:false"`
	SkipReason            *SkipReason `json:"skip_reason" gorm:"size:50"`
	PlayedAt              int64       `json:"played_at" gorm:"autoCreateTime;index"`
}
