package model

import "music-app-backend/pkg/model"

type Mood struct {
	model.BaseModel
	Name        string `json:"name" gorm:"unique;size:50"`
	Description string `json:"description"`
	ColorHex    string `json:"color_hex" gorm:"size:7"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}
