package model

import (
	"music-app-backend/pkg/model"
)

type PlaylistType string

const (
	PlaylistTypeCurated     PlaylistType = "curated"
	PlaylistTypeUserCreated PlaylistType = "user_created"
	PlaylistTypeMoodBased   PlaylistType = "mood_based"
)

type PlayList struct {
	model.BaseModel
	Name                 string       `json:"name" gorm:"not null;size:150"`
	Description          string       `json:"description"`
	ArtworkURL           string       `json:"artwork_url"`
	PlaylistType         PlaylistType `json:"playlist_type" gorm:"not null;size:50"`
	MoodID               *uint64      `json:"mood_id"`
	CreatedByUserID      *uint64      `json:"created_by_user_id"`
	IsPublic             bool         `json:"is_public" gorm:"default:true"`
	PlayCount            int64        `json:"play_count" gorm:"default:0"`
	SongCount            int          `json:"song_count" gorm:"default:0"`
	TotalDurationSeconds int          `json:"total_duration_seconds" gorm:"default:0"`
}
