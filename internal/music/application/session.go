package application

import (
	"context"
	model "music-app-backend/internal/music/domain"

	"gorm.io/gorm"
)

func (s *MusicService) CreateUploadSession(ctx context.Context, request *model.UploadSession) (*model.UploadSession, error) {
	err := s.repository.CreateUploadSession(ctx, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *MusicService) GetUploadSession(ctx context.Context, uploadID string) (*model.UploadSession, error) {
	uploadSession, err := s.repository.GetUploadSession(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	if uploadSession == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return uploadSession, nil
}

func (s *MusicService) UpdateUploadSession(ctx context.Context, uploadID string, updateData string) error {
	if err := s.repository.UpdateUploadSession(ctx, uploadID, updateData); err != nil {
		return err
	}
	return nil
}
