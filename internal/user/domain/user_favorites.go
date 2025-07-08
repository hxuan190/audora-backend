package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type UserFavorites struct {
	model.BaseModel
	UserID uuid.UUID `json:"user_id" gorm:"not null;index"`
	SongID uuid.UUID `json:"song_id" gorm:"not null;index"`
}