package model

import "music-app-backend/pkg/model"

type Genre struct {
	model.BaseModel
	Name        string `json:"name" gorm:"unique;size:50"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}
