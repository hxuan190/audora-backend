package repository

import (
	model "music-app-backend/internal/user/domain"
	"music-app-backend/pkg/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	database.Repository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.AutoMigrate(&model.User{}, &model.UserFavorites{}, &model.UserPreference{})
	return &UserRepository{
		Repository: database.NewRepository(db),
	}
}

func (r *UserRepository) IMockRepository() {}
