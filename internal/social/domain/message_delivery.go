package model

import (
	"music-app-backend/pkg/model"
)

type MessageDelivery struct {
	model.BaseModel
	MessageID   uint64 `json:"message_id" gorm:"not null;uniqueIndex:idx_message_user"`
	UserID      uint64 `json:"user_id" gorm:"not null;uniqueIndex:idx_message_user"`
	DeliveredAt int64  `json:"delivered_at"`
	ReadAt      *int64 `json:"read_at"`
}
