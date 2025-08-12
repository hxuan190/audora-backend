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
	UpdateSongProcessingResult(ctx context.Context, songID uint64, updates map[string]interface{}) error
	GetSongByID(ctx context.Context, songID uint64) (*model.Song, error)
	CreateProcessedAudioFormats(ctx context.Context, formats []model.ProcessedAudioFormat) error
	CreateAudioAnalysis(ctx context.Context, analysis *model.AudioAnalysis) error
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

func (db *MusicRepository) UpdateSongProcessingResult(ctx context.Context, songID uint64, updates map[string]interface{}) error {
	return db.db.WithContext(ctx).Model(&model.Song{}).Where("id = ?", songID).Updates(updates).Error
}

func (db *MusicRepository) GetSongByID(ctx context.Context, songID uint64) (*model.Song, error) {
	var song model.Song
	err := db.db.WithContext(ctx).Where("id = ?", songID).First(&song).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &song, nil
}

func (db *MusicRepository) CreateProcessedAudioFormats(ctx context.Context, formats []model.ProcessedAudioFormat) error {
	if len(formats) == 0 {
		return nil
	}
	return db.db.WithContext(ctx).Create(&formats).Error
}

func (db *MusicRepository) CreateAudioAnalysis(ctx context.Context, analysis *model.AudioAnalysis) error {
	return db.db.WithContext(ctx).Create(analysis).Error
}
