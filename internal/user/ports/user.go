package ports

type IUserRepository interface {
	IMockRepository()
}

type IUserService interface {
	IMockService()
}
