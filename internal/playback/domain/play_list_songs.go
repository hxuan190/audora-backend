package model

import (
	"music-app-backend/pkg/model"
)

type PlayListSongs struct {
	model.BaseModel
	PlaylistID    uint64  `json:"playlist_id" gorm:"not null;uniqueIndex:idx_playlist_song"`
	SongID        uint64  `json:"song_id" gorm:"not null;uniqueIndex:idx_playlist_song"`
	Position      int     `json:"position" gorm:"not null;uniqueIndex:idx_playlist_position"`
	AddedByUserID *uint64 `json:"added_by_user_id"`
}
