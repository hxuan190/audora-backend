package repository

import (
	"music-app-backend/pkg/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	database.Repository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: database.NewRepository(db),
	}
}

func (r *UserRepository) IMockRepository() {}
