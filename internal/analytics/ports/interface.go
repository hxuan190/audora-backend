package ports

type IAnalyticsRepository interface {
	IMockRepository()
}

type IAnalyticsService interface {
	IMockService()
}
