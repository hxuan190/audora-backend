package application

import "music-app-backend/internal/social/ports"

type SocialService struct {
	repository ports.ISocialRepository
}

func NewSocialService(repository ports.ISocialRepository) *SocialService {
	return &SocialService{repository: repository}
}

func (s *SocialService) IMockService() {}
