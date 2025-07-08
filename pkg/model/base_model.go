package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Status    string    `json:"status"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
