package model

import (
	"music-app-backend/pkg/model"

	"github.com/google/uuid"
)

type ContentTier string

const (
	ContentTierPublicDiscovery ContentTier = "public_discovery"
	ContentTierPremium         ContentTier = "premium"
	ContentTierExclusive       ContentTier = "exclusive"
)

type ProcessingStatus string

const (
	ProcessingStatusPending    ProcessingStatus = "pending"
	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusCompleted  ProcessingStatus = "completed"
	ProcessingStatusFailed     ProcessingStatus = "failed"
)

type Song struct {
	model.BaseModel
	ArtistID             uuid.UUID        `json:"artist_id" gorm:"not null"`
	Title                string           `json:"title" gorm:"not null;size:200"`
	Description          string           `json:"description"`
	FileURL              string           `json:"file_url" gorm:"not null"`
	FileSizeBytes        int64            `json:"file_size_bytes"`
	Duration             int64            `json:"duration_seconds"`
	ArtworkURL           string           `json:"artwork_url"`
	GenreID              int              `json:"genre_id"`
	MoodID               int              `json:"mood_id"`
	Tier                 ContentTier      `json:"tier" gorm:"not null;default:'public_discovery'"`
	AISuggestedTier      *ContentTier     `json:"ai_suggested_tier"`
	TierOverrideByArtist bool             `json:"tier_override_by_artist" gorm:"default:false"`
	BPM                  *int             `json:"bpm"`
	KeySignature         string           `json:"key_signature" gorm:"size:10"`
	IsExplicit           bool             `json:"is_explicit" gorm:"default:false"`
	IsProcessed          bool             `json:"is_processed" gorm:"default:false"`
	ProcessingStatus     ProcessingStatus `json:"processing_status" gorm:"default:'pending';size:50"`
	ProcessingError      string           `json:"processing_error"`
	PlayCount            int64            `json:"play_count" gorm:"default:0"`
	LikeCount            int              `json:"like_count" gorm:"default:0"`
	TipCount             int              `json:"tip_count" gorm:"default:0"`
	TotalTips            float64          `json:"total_tips" gorm:"type:decimal(10,2);default:0.00"`
	ReleaseDate          int64            `json:"release_date"`
	IsActive             bool             `json:"is_active" gorm:"default:true"`
}
