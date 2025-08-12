package repository

import (
	"context"
	model "music-app-backend/internal/music/domain"

	"gorm.io/gorm"
)

type IMusicRepository interface {
	InsertArtist(artist *model.Artist) error
	CreateUploadSession(ctx context.Context, upload *model.UploadSession) error
	GetUploadSession(ctx context.Context, uploadID string) (*model.UploadSession, error)
	UpdateUploadSession(ctx context.Context, uploadID string, updateData string) error
	InsertSong(song *model.Song) error
}

type MusicRepository struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) *MusicRepository {
	return &MusicRepository{
		db: db,
	}
}

func (db *MusicRepository) InsertArtist(artist *model.Artist) error {
	return db.db.Create(artist).Error
}

func (db *MusicRepository) CreateUploadSession(ctx context.Context, upload *model.UploadSession) error {
	return db.db.WithContext(ctx).Create(upload).Error
}

func (db *MusicRepository) GetUploadSession(ctx context.Context, uploadID string) (*model.UploadSession, error) {
	var uploadSession model.UploadSession
	err := db.db.WithContext(ctx).Where("id = ?", uploadID).First(&uploadSession).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &uploadSession, nil
}

func (db *MusicRepository) InsertSong(song *model.Song) error {
	return db.db.Create(song).Error
}

func (db *MusicRepository) UpdateUploadSession(ctx context.Context, uploadID string, updateData string) error {
	return db.db.WithContext(ctx).Model(&model.UploadSession{}).Where("id = ?", uploadID).Update("update_data", updateData).Error
}
