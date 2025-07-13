package application

import "music-app-backend/internal/playback/adapters/repository"

type PlaybackService struct {
	repository *repository.PlaybackRepository
}

func NewPlaybackService(repository *repository.PlaybackRepository) *PlaybackService {
	return &PlaybackService{repository: repository}
}

func (s *PlaybackService) IMockService() {}