package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type User struct {
	model.BaseModel
	KratosIdentityID uuid.UUID `json:"kratos_identity_id" gorm:"not null;unique;type:uuid"`
	Email            string    `json:"email" gorm:"not null;unique;size:255"`
	UserType         string    `json:"user_type" gorm:"not null;size:20;check:user_type IN ('artist', 'listener', 'admin')"`
	DisplayName      string    `json:"display_name" gorm:"size:100"`
	AvatarURL        string    `json:"avatar_url"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	LastLoginAt      *int64    `json:"last_login_at"`
}
