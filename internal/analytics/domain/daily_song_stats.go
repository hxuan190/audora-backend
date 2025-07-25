package model

import (
	"music-app-backend/pkg/model"
)

type DailySongStats struct {
	model.BaseModel
	SongID            uint64  `json:"song_id" gorm:"not null;uniqueIndex:idx_song_date"`
	Date              string  `json:"date" gorm:"not null;type:date;uniqueIndex:idx_song_date"`
	PlayCount         int     `json:"play_count" gorm:"default:0"`
	UniqueListeners   int     `json:"unique_listeners" gorm:"default:0"`
	CompletionRate    float64 `json:"completion_rate" gorm:"type:decimal(5,2);default:0.00"`
	AvgDurationPlayed int     `json:"avg_duration_played" gorm:"default:0"`
	SkipRate          float64 `json:"skip_rate" gorm:"type:decimal(5,2);default:0.00"`
}
