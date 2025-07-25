package model

import (
	"music-app-backend/pkg/model"
)

type UserFavorites struct {
	model.BaseModel
	UserID uint64 `json:"user_id" gorm:"not null;uniqueIndex:idx_user_song"`
	SongID uint64 `json:"song_id" gorm:"not null;uniqueIndex:idx_user_song"`
}
