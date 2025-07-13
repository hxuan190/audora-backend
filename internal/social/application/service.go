package application

import "music-app-backend/internal/social/adapters/repository"

type SocialService struct {
	repository *repository.SocialRepository
}

func NewSocialService(repository *repository.SocialRepository) *SocialService {
	return &SocialService{repository: repository}
}

func (s *SocialService) IMockService() {}
