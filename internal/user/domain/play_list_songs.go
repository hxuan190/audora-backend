package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type PlayListSongs struct {
	model.BaseModel
	PlaylistID    uuid.UUID  `json:"playlist_id" gorm:"not null;uniqueIndex:idx_playlist_song"`
	SongID        uuid.UUID  `json:"song_id" gorm:"not null;uniqueIndex:idx_playlist_song"`
	Position      int        `json:"position" gorm:"not null;uniqueIndex:idx_playlist_position"`
	AddedByUserID *uuid.UUID `json:"added_by_user_id"`
}
