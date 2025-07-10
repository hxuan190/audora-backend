package application

import "music-app-backend/internal/music/ports"

type MusicService struct {
	repository ports.IMusicRepository
}

func NewMusicService(repository ports.IMusicRepository) *MusicService {
	return &MusicService{repository: repository}
}

func (s *MusicService) IMockService() {}