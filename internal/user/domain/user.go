package model

import "music-app-backend/pkg/model"

type User struct {
	model.BaseModel
	KratosIdentityID string `json:"kratos_identity_id" gorm:"unique"`
	Email            string `json:"email" gorm:"unique"`
	UserType         string `json:"user_type"`
	DisplayName      string `json:"display_name"`
	AvatarURL        string `json:"avatar_url"`
	IsActive         bool   `json:"is_active" gorm:"default:true"`
	LastLoginAt      int64  `json:"last_login_at" gorm:"autoUpdateTime"`
}
