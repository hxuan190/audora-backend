package application

import "music-app-backend/internal/analytics/ports"

type AnalyticsService struct {
	repository ports.IAnalyticsRepository
}

func NewAnalyticsService(repository ports.IAnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repository: repository}
}

func (s *AnalyticsService) IMockService() {}
