package application

import "music-app-backend/internal/playback/ports"

type PlaybackService struct {
	repository ports.IPlaybackRepository
}

func NewPlaybackService(repository ports.IPlaybackRepository) *PlaybackService {
	return &PlaybackService{repository: repository}
}

func (s *PlaybackService) IMockService() {}