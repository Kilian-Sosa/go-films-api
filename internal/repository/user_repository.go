package repository

import (
	"errors"

	"gorm.io/gorm"

	"go-films-api/internal/domain"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
}

type userRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) UserRepository {
	return &userRepositoryGorm{db: db}
}

func (r *userRepositoryGorm) CreateUser(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryGorm) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
