package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type TargetType string

const (
	TargetTypeAllActiveListeners    TargetType = "all_active_listeners"
	TargetTypeSpecificSongListeners TargetType = "specific_song_listeners"
	TargetTypeFollowers             TargetType = "followers"
)

type ArtistMessage struct {
	model.BaseModel
	ArtistID     uuid.UUID  `json:"artist_id" gorm:"type:uuid;not null"`
	MessageText  string     `json:"message_text" gorm:"not null"`
	TargetType   TargetType `json:"target_type" gorm:"not null;size:50"`
	TargetSongID *uuid.UUID `json:"target_song_id" gorm:"type:uuid"`
	SentToCount  int        `json:"sent_to_count" gorm:"default:0"`
	ReadCount    int        `json:"read_count" gorm:"default:0"`
}
