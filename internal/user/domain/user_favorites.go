package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type UserFavorites struct {
	model.BaseModel
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_song"`
	SongID uuid.UUID `json:"song_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_song"`
}
