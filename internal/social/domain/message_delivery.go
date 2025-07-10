package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type MessageDelivery struct {
	model.BaseModel
	MessageID   uuid.UUID `json:"message_id" gorm:"not null;uniqueIndex:idx_message_user"`
	UserID      uuid.UUID `json:"user_id" gorm:"not null;uniqueIndex:idx_message_user"`
	DeliveredAt int64     `json:"delivered_at"`
	ReadAt      *int64    `json:"read_at"`
}
