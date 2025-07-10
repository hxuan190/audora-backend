package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;"`
	Status    string    `json:"status"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewBaseModel() *BaseModel {
	uuid, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("Error creating base model: %v", err)
	}

	return &BaseModel{
		ID:        uuid,
		Status:    "active",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
