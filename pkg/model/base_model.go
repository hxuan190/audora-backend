package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64 `json:"id" gorm:"primaryKey;"`
	Status    string `json:"status" gorm:"default:active;check:status IN ('active', 'inactive')"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewBaseModel() *BaseModel {
	uuid := uint64(time.Now().Unix())

	return &BaseModel{
		ID:        uuid,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uint64(time.Now().Unix())
	return
}
