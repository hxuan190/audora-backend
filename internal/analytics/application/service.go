package application

import "music-app-backend/internal/analytics/adapters/repository"

type AnalyticsService struct {
	repository *repository.AnalyticsRepository
}

func NewAnalyticsService(repository *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repository: repository}
}	

func (s *AnalyticsService) IMockService() {}
