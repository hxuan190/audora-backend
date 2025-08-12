package model

import "time"

type UploadSession struct {
	ID         string    `json:"id" gorm:"primaryKey;"`
	ArtistID   uint64    `json:"artist_id" gorm:"not null;"`
	UserID     uint64    `json:"user_id" gorm:"not null;"`
	Filename   string    `json:"filename" gorm:"not null;"`
	FileSize   int64     `json:"file_size" gorm:"not null;"`
	ObjectPath string    `json:"object_path" gorm:"not null;"`
	Status     string    `json:"status" gorm:"not null;"`
	ExpiresAt  time.Time `json:"expires_at" gorm:"not null;"`
}
