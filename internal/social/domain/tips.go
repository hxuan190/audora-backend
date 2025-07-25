package model

import (
	"music-app-backend/pkg/model"
)

type TipStatus string

const (
	TipStatusPending   TipStatus = "pending"
	TipStatusCompleted TipStatus = "completed"
	TipStatusFailed    TipStatus = "failed"
	TipStatusRefunded  TipStatus = "refunded"
)

type Tip struct {
	model.BaseModel
	FromUserID            uint64    `json:"from_user_id" gorm:"not null"`
	ToArtistID            uint64    `json:"to_artist_id" gorm:"not null"`
	SongID                *uint64   `json:"song_id"`
	AmountCents           int       `json:"amount_cents" gorm:"not null"`
	Currency              string    `json:"currency" gorm:"default:'USD';size:3"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id" gorm:"unique;size:100"`
	StripeChargeID        string    `json:"stripe_charge_id" gorm:"size:100"`
	PlatformFeeCents      int       `json:"platform_fee_cents" gorm:"not null"`
	ArtistPayoutCents     int       `json:"artist_payout_cents" gorm:"not null"`
	Status                TipStatus `json:"status" gorm:"default:'pending';size:50"`
	Message               string    `json:"message"`
	IsAnonymous           bool      `json:"is_anonymous" gorm:"default:false"`
	ProcessedAt           *int64    `json:"processed_at"`
}
