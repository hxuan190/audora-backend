package application

import "music-app-backend/internal/music/adapters/repository"

type MusicService struct {
	repository *repository.MusicRepository
}

func NewMusicService(repository *repository.MusicRepository) *MusicService {
	return &MusicService{repository: repository}
}

func (s *MusicService) IMockService() {}