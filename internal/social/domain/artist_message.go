package model

import (
	"music-app-backend/pkg/model"
)

type TargetType string

const (
	TargetTypeAllActiveListeners    TargetType = "all_active_listeners"
	TargetTypeSpecificSongListeners TargetType = "specific_song_listeners"
	TargetTypeFollowers             TargetType = "followers"
)

type ArtistMessage struct {
	model.BaseModel
	ArtistID     uint64     `json:"artist_id" gorm:"not null"`
	MessageText  string     `json:"message_text" gorm:"not null"`
	TargetType   TargetType `json:"target_type" gorm:"not null;size:50"`
	TargetSongID *uint64    `json:"target_song_id"`
	SentToCount  int        `json:"sent_to_count" gorm:"default:0"`
	ReadCount    int        `json:"read_count" gorm:"default:0"`
}
